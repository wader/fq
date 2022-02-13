def tobits: _tobits(1; false; 0);
def tobytes: _tobits(8; false; 0);
def tobitsrange: _tobits(1; true; 0);
def tobytesrange: _tobits(8; true; 0);
def tobits($pad): _tobits(1; false; $pad);
def tobytes($pad): _tobits(8; false; $pad);
def tobitsrange($pad): _tobits(1; true; $pad);
def tobytesrange($pad): _tobits(8; true; $pad);