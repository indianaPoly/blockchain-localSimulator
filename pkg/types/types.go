package types

type Block struct {
    Index        int
    Timestamp    string
    Data         string
    PreviousHash string
}

type Blockchain struct {
    Blocks []Block
}