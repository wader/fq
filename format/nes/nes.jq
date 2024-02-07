def nes_toasm:
  select(.op_code and (.args // "")) | "\(.op_code) \(.args // "")";
def nes_tokitty($size):
  def _sub:
    ( gsub("00"; "\u0022\u0022\u0022")
    | gsub("01"; "\u007F\u0022\u0022")
    | gsub("10"; "\u0022\u007F\u0022")
    | gsub("11"; "\u0022\u0022\u007F")
    | gsub(" "; "")
    );
  ( "\(.combined | _sub)"
  | "\\e_Ga=T,f=24,s=8,v=8,c=\(2*$size),r=\(1*$size);" + @base64 + "\\e\\\\\\n"
  );
def nes_tokitty:
  nes_tokitty(1);
