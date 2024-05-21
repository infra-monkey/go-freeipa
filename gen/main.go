// Copyright © 2022 IN2P3 Computing Centre, IN2P3, CNRS
// Copyright © 2018 Philippe Voinov
//
// Contributor(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2021
//
// This software is governed by the CeCILL license under French law and
// abiding by the rules of distribution of free software.  You can  use,
// modify and/ or redistribute the software under the terms of the CeCILL
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info".
//
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability.
//
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or
// data to be ensured and,  more generally, to use and operate it in the
// same conditions as regards security.
//
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL license and that you accept its terms.

package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var skipCommands = []string{
	"batch",
	"command_defaults",
	"env",
}

var skipClasses = []string{
	"class",
	"command",
}

// Only assume {"type": "dict"} to be the receiver (eg. "user" in "user_add")
// for the following commands.
var dictResultCommands = []string{
	"add",
	"find",
	"mod",
	"show",
}

func main() {
	if e := actualMain(); e != nil {
		log.Fatalf("%v", e)
	}
}

func actualMain() error {
	overrides, err := loadDirtyOverrides()
	if err != nil {
		return errors.Wrap(err, "loading local overrides")
	}

	schema, e := loadSchema(overrides)
	if e != nil {
		return e
	}
	errs, e := loadErrs()
	if e != nil {
		return e
	}
	if e := generateMain(schema, errs); e != nil {
		return e
	}
	return nil
}

func loadSchema(localOverrides DirtyOverrides) (*Schema, error) {
	input, e := os.ReadFile("../data/schema.json")
	if e != nil {
		return nil, e
	}
	parsed := &SchemaDump{}
	if e = json.Unmarshal(input, parsed); e != nil {
		return nil, e
	}
	schema := parsed.Result.Result

	cmds := make([]*Command, 0, len(schema.Commands))
	for _, c := range schema.Commands {
		var skip bool
		for _, s := range skipCommands {
			if s == c.Name {
				skip = true
			}
		}
		if !skip {
			// HACK Many result values for FreeIPA methods
			// have {"type": "dict"}. Often this means the type
			// of the receiver (eg. "user" in "user_add"), but not
			// always. Limit this guessing to whitelisted method types.
			var guessDictRes bool
			for _, v := range dictResultCommands {
				if c.AttrName == v {
					guessDictRes = true
				}
			}
			if guessDictRes {
				for _, p := range c.Output {
					if p.Name == "result" && p.Type == "dict" {
						p.Type = "dict_guess_receiver"
					}
				}
			}

			cmds = append(cmds, c)
		}
	}
	schema.Commands = cmds

	classes := make([]*Class, 0, len(schema.Classes))
	for _, c := range schema.Classes {
		var skip bool
		for _, s := range skipCommands {
			if s == c.Name {
				skip = true
			}
		}

		// HACK !!
		// do we have local maintained overrides for this particular class ?
		var classOverrides *ClassOverrides
		classOverridesRaw, ok := localOverrides.Classes[c.Name]
		if ok {
			classOverrides = &classOverridesRaw
		}

		if !skip {
			// HACK FreeIPA admin user has no "givenname" or "cn", even though the schema
			// says these fields are required. This workaround makes it optional.
			if c.Name == "user" {
				for _, p := range c.Params {
					if p.Name == "givenname" || p.Name == "cn" {
						v := false
						p.RequiredRaw = &v
					}
				}
			}

			// HACK FreeIPA host has several fields which are not required.
			hostNotRequiredParams := []string{
				"subject",
				"serial_number",
				"serial_number_hex",
				"issuer",
				"valid_not_before",
				"valid_not_after",
				"md5_fingerprint",
				"sha1_fingerprint",
				"sha256_fingerprint",
				"managing_host",
				"ipaallowedtoperform_read_keys_user",
				"ipaallowedtoperform_read_keys_group",
				"ipaallowedtoperform_read_keys_host",
				"ipaallowedtoperform_read_keys_hostgroup",
				"ipaallowedtoperform_write_keys_user",
				"ipaallowedtoperform_write_keys_group",
				"ipaallowedtoperform_write_keys_host",
				"ipaallowedtoperform_write_keys_hostgroup",
			}
			if c.Name == "host" {
				for _, p := range c.Params {
					for _, pp := range hostNotRequiredParams {
						if p.Name == pp {
							v := false
							p.RequiredRaw = &v
						}
					}
				}
			}

			// HACK FreeIPA sometimes doesn't supply boolean fields which are
			// marked required in schema. This workaround makes them optional.
			for _, p := range c.Params {
				if p.Type == "bool" {
					v := false
					p.RequiredRaw = &v
				}
			}

			// HACK Fields starting with "member_" or "memberof_" generally seem to be multivalued,
			// even though the schema doesn't say so. Assuming they are multivalued
			// will work even if they end up actually being single-valued.
			for _, p := range c.Params {
				multiValuePrefixes := []string{
					"member_",
					"memberof_",
					"memberindirect_",
					"memberofindirect_",
					"memberuser_",
					"memberhost_",
					"membercmd_",
					"memberallowcmd_",
					"memberdenycmd_",
				}
				for _, multiValuePrefix := range multiValuePrefixes {
					if strings.HasPrefix(p.Name, multiValuePrefix) {
						p.Multivalue = true
					}
				}
			}

			// HACK: Use our local overrides
			if classOverrides != nil {
				// HACK use local overrides to modify directly the params
				for _, p := range c.Params {
					paramOverrides, ok := classOverrides.Params[p.Name]
					if !ok {
						continue
					}

					paramOverrides.OverrideParams(p)
				}
			}

			classes = append(classes, c)
		}
	}

	// TODO assert "version" on all objects is "1"
	// TODO assert that names are consistent within each object

	return &schema, nil
}

func loadErrs() ([]ErrDesc, error) {
	in, e := os.ReadFile("../data/errors.json")
	if e != nil {
		return nil, e
	}
	var out []ErrDesc
	if e = json.Unmarshal(in, &out); e != nil {
		return nil, e
	}
	return out, nil
}

func generateMain(schema *Schema, errs []ErrDesc) error {
	t, e := template.New("freeipa.gotmpl").Funcs(template.FuncMap{
		"ToGoType":  toGoType,
		"TrimSpace": strings.TrimSpace,
	}).ParseFiles("./freeipa.gotmpl")
	if e != nil {
		return e
	}
	f, e := os.OpenFile("../freeipa/generated.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if e != nil {
		return e
	}
	e = t.Execute(f, struct {
		Schema *Schema
		Errs   []ErrDesc
	}{schema, errs})
	if e != nil {
		return e
	}
	return nil
}

func stringSliceContains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}
