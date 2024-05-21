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
	"os"
)

type DirtyOverrides struct {
	Classes map[string]ClassOverrides `json:"classes"`
}

type ClassOverrides struct {
	Params map[string]ClassParamsOverrides `json:"params"`
}

type ClassParamsOverrides struct {
	Required   *bool `json:"required,omitempty"`
	Multivalue *bool `json:"multivalue,omitempty"`
}

func (c ClassParamsOverrides) OverrideParams(p *Param) {
	if c.Required != nil {
		p.RequiredRaw = c.Required
	}

	if c.Multivalue != nil {
		p.Multivalue = *c.Multivalue
	}
}

func loadDirtyOverrides() (DirtyOverrides, error) {
	var overrides DirtyOverrides
	in, e := os.ReadFile("../data/dirty_overrides.json")
	if e != nil {
		return overrides, e
	}
	if e = json.Unmarshal(in, &overrides); e != nil {
		return overrides, e
	}
	return overrides, nil
}
