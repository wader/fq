import torch
from safetensors.torch import save_file

tensors = {
    "weight1": torch.reshape(torch.arange(12, dtype=torch.float32), (12,)),
    "weight2": torch.reshape(torch.arange(12, dtype=torch.int64), (3, 4)),
    "weight3": torch.reshape(torch.arange(12, dtype=torch.float16), (2, 2, 3)),
    "weight4": torch.reshape(torch.arange(12, dtype=torch.bfloat16), (4, 3)),
}

save_file(tensors, "format/safetensors/testdata/test.safetensors")