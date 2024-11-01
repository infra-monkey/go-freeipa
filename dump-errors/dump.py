# Copyright © 2022 IN2P3 Computing Centre, IN2P3, CNRS
# Copyright © 2018 Philippe Voinov
#
# Contributor(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2021
#
# This software is governed by the CeCILL license under French law and
# abiding by the rules of distribution of free software.  You can  use,
# modify and/ or redistribute the software under the terms of the CeCILL
# license as circulated by CEA, CNRS and INRIA at the following URL
# "http://www.cecill.info".
#
# As a counterpart to the access to the source code and  rights to copy,
# modify and redistribute granted by the license, users are provided only
# with a limited warranty  and the software's author,  the holder of the
# economic rights,  and the successive licensors  have only  limited
# liability.
#
# In this respect, the user's attention is drawn to the risks associated
# with loading,  using,  modifying and/or developing or reproducing the
# software by the user in light of its specific status of free software,
# that may mean  that it is complicated to manipulate,  and  that  also
# therefore means  that it is reserved for developers  and  experienced
# professionals having in-depth computer knowledge. Users are therefore
# encouraged to load and test the software's suitability as regards their
# requirements in conditions enabling the security of their systems and/or
# data to be ensured and,  more generally, to use and operate it in the
# same conditions as regards security.
# 
# The fact that you are presently reading this means that you have had
# knowledge of the CeCILL license and that you accept its terms.

import urllib.request
import sys
import imp
import re
import inspect
import json

ERRORS_PY_URL = (
    "https://raw.githubusercontent.com/freeipa/freeipa/ipa-4-12/ipalib/errors.py"
)

import_regex = re.compile(r"^(from [\w\.]+ )?import \w+( as \w+)?$")


def should_keep(l):
    return import_regex.match(l) is None


errors_py_str = urllib.request.urlopen(ERRORS_PY_URL).read().decode("utf-8")
errors_py_str = "\n".join([l for l in errors_py_str.splitlines() if should_keep(l)])
errors_py_str = (
    """
class Six:
    PY3 = True
six = Six()
ungettext = None
class Messages:
    def iter_messages(*args):
        return []
messages = Messages()
"""
    + errors_py_str
)

errors_mod = imp.new_module("errors")
exec(errors_py_str, errors_mod.__dict__)

error_codes = [
    {"name": k, "errno": v.errno}
    for k, v in inspect.getmembers(errors_mod)
    if hasattr(v, "__dict__") and type(v.__dict__.get("errno", None)) == int
]
error_codes.sort(key=lambda x: x["errno"])

with open("../data/errors.json", "w") as f:
    json.dump(error_codes, f)
