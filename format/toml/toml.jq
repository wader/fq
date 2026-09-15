def _toml__todisplay: tovalue;
# 2 is default for BurntSushi/toml
def to_toml($opts): _to_toml({"indent": 2} + $opts);
def to_toml: to_toml({});
