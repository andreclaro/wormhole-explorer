package guardiansets

import (
	"errors"
	"fmt"
	"time"

	"github.com/certusone/wormhole/node/pkg/common"
	sdk "github.com/wormhole-foundation/wormhole/sdk/vaa"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/wormhole-foundation/wormhole-explorer/fly/config"
)

// GuardianSetHistory contains information about all guardian sets for the current network (past and present).
type GuardianSetHistory struct {
	guardianSetsByIndex    []common.GuardianSet
	expirationTimesByIndex []time.Time
}

// Verify takes a VAA as input and validates its guardian signatures.
func (h *GuardianSetHistory) Verify(vaa *sdk.VAA) error {

	idx := vaa.GuardianSetIndex

	// Make sure the index exists
	if idx >= uint32(len(h.guardianSetsByIndex)) {
		return fmt.Errorf("Guardian Set Index is out of bounds: got %d, max is %d",
			vaa.GuardianSetIndex,
			len(h.guardianSetsByIndex),
		)
	}

	// Verify guardian signatures
	if sdk.VerifySignatures(vaa.SigningMsg().Bytes(), vaa.Signatures, h.guardianSetsByIndex[idx].Keys) {
		return nil
	} else {
		return errors.New("VAA contains invalid signatures")
	}
}

// GetLatest returns the lastest guardian set.
func (h GuardianSetHistory) GetLatest() common.GuardianSet {
	return h.guardianSetsByIndex[len(h.guardianSetsByIndex)-1]
}

// Get get guardianset config by enviroment.
func GetByEnv(enviroment string) GuardianSetHistory {
	switch enviroment {
	case config.P2pTestNet:
		return getTestnetGuardianSet()
	default:
		return getMainnetGuardianSet()
	}
}

func getTestnetGuardianSet() GuardianSetHistory {
	const tenYears = time.Hour * 24 * 365 * 10
	gs0TestValidUntil := time.Now().Add(tenYears)
	gstest0 := common.GuardianSet{
		Index: 0,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x13947Bd48b18E53fdAeEe77F3473391aC727C638"), //
		},
	}
	return GuardianSetHistory{
		guardianSetsByIndex:    []common.GuardianSet{gstest0},
		expirationTimesByIndex: []time.Time{gs0TestValidUntil},
	}
}

func getMainnetGuardianSet() GuardianSetHistory {
	gs0ValidUntil := time.Unix(1628599904, 0) // Tue Aug 10 2021 12:51:44 GMT+0000
	gs0 := common.GuardianSet{
		Index: 0,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"), // Certus One
		},
	}

	gs1ValidUntil := time.Unix(1650566103, 0) // Thu Apr 21 2022 18:35:03 GMT+0000
	gs1 := common.GuardianSet{
		Index: 1,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"), // Certus One
			eth_common.HexToAddress("0xfF6CB952589BDE862c25Ef4392132fb9D4A42157"), // Staked
			eth_common.HexToAddress("0x114De8460193bdf3A2fCf81f86a09765F4762fD1"), // Figment
			eth_common.HexToAddress("0x107A0086b32d7A0977926A205131d8731D39cbEB"), // ChainodeTech
			eth_common.HexToAddress("0x8C82B2fd82FaeD2711d59AF0F2499D16e726f6b2"), // Inotel
			eth_common.HexToAddress("0x11b39756C042441BE6D8650b69b54EbE715E2343"), // HashQuark
			eth_common.HexToAddress("0x54Ce5B4D348fb74B958e8966e2ec3dBd4958a7cd"), // ChainLayer
			eth_common.HexToAddress("0xeB5F7389Fa26941519f0863349C223b73a6DDEE7"), // DokiaCapital
			eth_common.HexToAddress("0x74a3bf913953D695260D88BC1aA25A4eeE363ef0"), // Forbole
			eth_common.HexToAddress("0x000aC0076727b35FBea2dAc28fEE5cCB0fEA768e"), // Staking Fund
			eth_common.HexToAddress("0xAF45Ced136b9D9e24903464AE889F5C8a723FC14"), // MoonletWallet
			eth_common.HexToAddress("0xf93124b7c738843CBB89E864c862c38cddCccF95"), // P2P Validator
			eth_common.HexToAddress("0xD2CC37A4dc036a8D232b48f62cDD4731412f4890"), // 01node
			eth_common.HexToAddress("0xDA798F6896A3331F64b48c12D1D57Fd9cbe70811"), // MCF-V2-MAINNET
			eth_common.HexToAddress("0x71AA1BE1D36CaFE3867910F99C09e347899C19C3"), // Everstake
			eth_common.HexToAddress("0x8192b6E7387CCd768277c17DAb1b7a5027c0b3Cf"), // Chorus One
			eth_common.HexToAddress("0x178e21ad2E77AE06711549CFBB1f9c7a9d8096e8"), // syncnode
			eth_common.HexToAddress("0x5E1487F35515d02A92753504a8D75471b9f49EdB"), // Triton
			eth_common.HexToAddress("0x6FbEBc898F403E4773E95feB15E80C9A99c8348d"), // Staking Facilities
		},
	}

	const tenYears = time.Hour * 24 * 365 * 10
	gs2ValidUntil := time.Now().Add(tenYears) // still valid so we add 10 years
	gs2 := common.GuardianSet{
		Index: 2,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"), // Certus One
			eth_common.HexToAddress("0xfF6CB952589BDE862c25Ef4392132fb9D4A42157"), // Staked
			eth_common.HexToAddress("0x114De8460193bdf3A2fCf81f86a09765F4762fD1"), // Figment
			eth_common.HexToAddress("0x107A0086b32d7A0977926A205131d8731D39cbEB"), // ChainodeTech
			eth_common.HexToAddress("0x8C82B2fd82FaeD2711d59AF0F2499D16e726f6b2"), // Inotel
			eth_common.HexToAddress("0x11b39756C042441BE6D8650b69b54EbE715E2343"), // HashQuark
			eth_common.HexToAddress("0x54Ce5B4D348fb74B958e8966e2ec3dBd4958a7cd"), // ChainLayer
			eth_common.HexToAddress("0x66B9590e1c41e0B226937bf9217D1d67Fd4E91F5"), // FTX
			eth_common.HexToAddress("0x74a3bf913953D695260D88BC1aA25A4eeE363ef0"), // Forbole
			eth_common.HexToAddress("0x000aC0076727b35FBea2dAc28fEE5cCB0fEA768e"), // Staking Fund
			eth_common.HexToAddress("0xAF45Ced136b9D9e24903464AE889F5C8a723FC14"), // MoonletWallet
			eth_common.HexToAddress("0xf93124b7c738843CBB89E864c862c38cddCccF95"), // P2P Validator
			eth_common.HexToAddress("0xD2CC37A4dc036a8D232b48f62cDD4731412f4890"), // 01node
			eth_common.HexToAddress("0xDA798F6896A3331F64b48c12D1D57Fd9cbe70811"), // MCF-V2-MAINNET
			eth_common.HexToAddress("0x71AA1BE1D36CaFE3867910F99C09e347899C19C3"), // Everstake
			eth_common.HexToAddress("0x8192b6E7387CCd768277c17DAb1b7a5027c0b3Cf"), // Chorus One
			eth_common.HexToAddress("0x178e21ad2E77AE06711549CFBB1f9c7a9d8096e8"), // syncnode
			eth_common.HexToAddress("0x5E1487F35515d02A92753504a8D75471b9f49EdB"), // Triton
			eth_common.HexToAddress("0x6FbEBc898F403E4773E95feB15E80C9A99c8348d"), // Staking Facilities
			// devnet
			// eth_common.HexToAddress("0xbeFA429d57cD18b7F8A4d91A2da9AB4AF05d0FBe"),
		},
	}

	gs3ValidUntil := time.Now().Add(tenYears) // still valid so we add 10 years
	gs3 := common.GuardianSet{
		Index: 3,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"), // Certus One
			eth_common.HexToAddress("0xfF6CB952589BDE862c25Ef4392132fb9D4A42157"), // Staked
			eth_common.HexToAddress("0x114De8460193bdf3A2fCf81f86a09765F4762fD1"), // Figment
			eth_common.HexToAddress("0x107A0086b32d7A0977926A205131d8731D39cbEB"), // ChainodeTech
			eth_common.HexToAddress("0x8C82B2fd82FaeD2711d59AF0F2499D16e726f6b2"), // Inotel
			eth_common.HexToAddress("0x11b39756C042441BE6D8650b69b54EbE715E2343"), // HashQuark
			eth_common.HexToAddress("0x54Ce5B4D348fb74B958e8966e2ec3dBd4958a7cd"), // ChainLayer
			eth_common.HexToAddress("0x15e7cAF07C4e3DC8e7C469f92C8Cd88FB8005a20"), // xLabs
			eth_common.HexToAddress("0x74a3bf913953D695260D88BC1aA25A4eeE363ef0"), // Forbole
			eth_common.HexToAddress("0x000aC0076727b35FBea2dAc28fEE5cCB0fEA768e"), // Staking Fund
			eth_common.HexToAddress("0xAF45Ced136b9D9e24903464AE889F5C8a723FC14"), // MoonletWallet
			eth_common.HexToAddress("0xf93124b7c738843CBB89E864c862c38cddCccF95"), // P2P Validator
			eth_common.HexToAddress("0xD2CC37A4dc036a8D232b48f62cDD4731412f4890"), // 01node
			eth_common.HexToAddress("0xDA798F6896A3331F64b48c12D1D57Fd9cbe70811"), // MCF-V2-MAINNET
			eth_common.HexToAddress("0x71AA1BE1D36CaFE3867910F99C09e347899C19C3"), // Everstake
			eth_common.HexToAddress("0x8192b6E7387CCd768277c17DAb1b7a5027c0b3Cf"), // Chorus One
			eth_common.HexToAddress("0x178e21ad2E77AE06711549CFBB1f9c7a9d8096e8"), // syncnode
			eth_common.HexToAddress("0x5E1487F35515d02A92753504a8D75471b9f49EdB"), // Triton
			eth_common.HexToAddress("0x6FbEBc898F403E4773E95feB15E80C9A99c8348d"), // Staking Facilities

		},
	}

	return GuardianSetHistory{
		guardianSetsByIndex:    []common.GuardianSet{gs0, gs1, gs2, gs3},
		expirationTimesByIndex: []time.Time{gs0ValidUntil, gs1ValidUntil, gs2ValidUntil, gs3ValidUntil},
	}
}
