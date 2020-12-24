package bridgesig

import (
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"math/big"

	"github.com/incognito-core-libs/go-ethereum/common/hexutil"
	ethcrypto "github.com/incognito-core-libs/go-ethereum/crypto"
)

func DecodeECDSASig(sigStr string) (
	v byte,
	r string,
	s string,
	err error,
) {
	sig, ver, errDecode := base58.Base58Check{}.Decode(sigStr)
	if (len(sig) != CBridgeSigSz) || (ver != common.ZeroByte) || (errDecode != nil) {
		err = NewBriSignatureError(InvalidInputParamsSizeErr, nil)
		return
	}
	v = byte(sig[64] + 27)
	r = hexutil.Encode(sig[:32])
	s = hexutil.Encode(sig[32:64])
	return
}

// B2ImN is Bytes to Int mod N, with N is secp256k1 curve order
func B2ImN(bytes []byte) *big.Int {
	x := big.NewInt(0)
	x.SetBytes(ethcrypto.Keccak256Hash(bytes).Bytes())
	for x.Cmp(ethcrypto.S256().Params().N) != -1 {
		x.SetBytes(ethcrypto.Keccak256Hash(x.Bytes()).Bytes())
	}
	return x
}
