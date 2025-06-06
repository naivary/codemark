# Usgae of available marker


// +path:to:marker=3 -> i[8,16,32,64]|*i[8,16,32,64] & uint[8,16,32,64]|*uint[8,16,32,64]
// +path:to:marker=0x23ef -> i[8,16,32,64]|*i[8,16,32,64] & uint[8,16,32,64]|*uint[8,16,32,64]
// +path:to:marker=0o352 -> i[8,16,32,64]|*i[8,16,32,64] & uint[8,16,32,64]|*uint[8,16,32,64]

// +path:to:marker=3.0 -> float[32,64]|*float[32,64]
// +path:to:marker=2+3i-> complex[64,128]|*complex[64,128]
