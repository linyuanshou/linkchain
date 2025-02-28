package ringct

import (
	"encoding/hex"
	"fmt"
	"github.com/lianxiangcloud/linkchain/libs/cryptonote/types"
	"github.com/lianxiangcloud/linkchain/libs/cryptonote/xcrypto"
	"testing"
)

var (
	message = types.Key{25, 15, 184, 52, 135, 221, 0, 36, 2, 123, 121, 179, 196, 164, 211, 38, 206, 125, 82, 57, 123, 84, 128, 216, 22, 67, 176, 231, 151, 51, 179, 58}
	mixring = types.CtkeyM{
		{
			{types.Key{191, 106, 61, 16, 27, 205, 154, 179, 20, 74, 153, 245, 60, 147, 226, 119, 91, 87, 96, 247, 21, 79, 251, 241, 188, 61, 206, 122, 110, 222, 72, 99,}, types.Key{224, 99, 57, 131, 42, 137, 71, 62, 118, 246, 63, 141, 83, 195, 40, 70, 234, 7, 145, 208, 182, 171, 44, 28, 231, 35, 153, 213, 163, 50, 207, 226}},
			{types.Key{139, 77, 189, 180, 37, 241, 187, 94, 61, 232, 71, 9, 170, 120, 192, 233, 1, 8, 189, 2, 56, 23, 195, 217, 134, 162, 255, 154, 80, 128, 88, 167,}, types.Key{173, 190, 4, 158, 54, 219, 41, 106, 233, 135, 113, 115, 236, 142, 133, 155, 78, 184, 197, 227, 200, 78, 39, 250, 250, 142, 111, 42, 79, 109, 239, 57}},
			{types.Key{210, 102, 211, 101, 247, 92, 184, 222, 202, 163, 13, 218, 88, 164, 123, 242, 89, 152, 133, 216, 214, 124, 137, 174, 3, 28, 23, 51, 87, 134, 154, 2,}, types.Key{195, 188, 86, 161, 255, 169, 204, 219, 75, 208, 130, 148, 28, 206, 17, 148, 156, 96, 64, 194, 25, 155, 228, 119, 103, 59, 247, 117, 177, 95, 51, 137}},
			{types.Key{11, 1, 173, 171, 65, 114, 209, 100, 196, 217, 56, 187, 165, 196, 209, 46, 115, 66, 104, 253, 134, 72, 47, 157, 185, 118, 220, 2, 217, 116, 215, 94,}, types.Key{56, 201, 82, 79, 207, 129, 22, 96, 118, 41, 219, 9, 103, 10, 3, 109, 47, 195, 235, 52, 84, 255, 74, 212, 203, 33, 162, 206, 219, 85, 229, 183}},
			{types.Key{49, 119, 211, 36, 87, 40, 124, 88, 200, 135, 252, 185, 2, 175, 111, 41, 38, 159, 216, 111, 28, 139, 105, 80, 9, 195, 91, 65, 50, 34, 56, 218,}, types.Key{61, 24, 113, 42, 22, 146, 152, 244, 127, 71, 149, 242, 184, 52, 147, 21, 219, 33, 198, 42, 45, 138, 238, 225, 86, 185, 44, 196, 76, 133, 163, 0}},
			{types.Key{3, 103, 53, 107, 59, 21, 167, 252, 164, 0, 212, 26, 176, 163, 76, 133, 23, 73, 94, 13, 34, 62, 4, 213, 50, 136, 39, 26, 237, 194, 10, 170,}, types.Key{6, 9, 18, 202, 0, 161, 249, 221, 23, 41, 6, 199, 236, 200, 177, 32, 206, 181, 186, 4, 184, 188, 13, 189, 16, 3, 120, 29, 209, 185, 208, 250}},
			{types.Key{224, 113, 124, 72, 210, 71, 82, 70, 134, 182, 40, 155, 163, 173, 131, 12, 119, 56, 166, 192, 239, 43, 144, 189, 223, 145, 132, 188, 242, 153, 62, 172,}, types.Key{17, 199, 75, 38, 219, 220, 31, 135, 137, 64, 76, 55, 109, 148, 98, 214, 49, 166, 92, 40, 33, 73, 119, 228, 125, 210, 85, 122, 52, 33, 56, 194}},
			{types.Key{47, 125, 54, 39, 93, 85, 193, 230, 197, 198, 129, 229, 95, 83, 71, 83, 23, 166, 89, 207, 125, 48, 235, 236, 61, 203, 30, 1, 230, 72, 200, 135,}, types.Key{46, 102, 164, 13, 15, 16, 81, 194, 150, 165, 158, 20, 118, 126, 192, 191, 143, 155, 162, 164, 254, 4, 108, 12, 234, 237, 30, 66, 170, 220, 240, 165}},
			{types.Key{28, 119, 123, 243, 165, 20, 4, 32, 229, 13, 217, 25, 83, 18, 68, 71, 181, 246, 48, 21, 194, 80, 83, 92, 31, 60, 234, 170, 210, 67, 1, 155,}, types.Key{86, 223, 17, 236, 251, 100, 5, 19, 218, 8, 185, 107, 171, 195, 142, 245, 217, 202, 8, 154, 214, 26, 174, 178, 26, 245, 42, 164, 152, 159, 64, 4}},
			{types.Key{120, 163, 243, 180, 139, 203, 12, 177, 184, 45, 252, 86, 181, 16, 154, 169, 171, 147, 145, 122, 162, 198, 231, 212, 200, 80, 143, 156, 218, 169, 232, 184,}, types.Key{197, 206, 204, 189, 222, 4, 169, 178, 164, 209, 255, 26, 135, 218, 183, 86, 19, 217, 179, 234, 177, 67, 28, 181, 38, 35, 193, 64, 74, 203, 140, 77}},
			{types.Key{220, 145, 117, 210, 42, 173, 189, 127, 220, 9, 197, 113, 151, 121, 124, 151, 230, 37, 244, 147, 44, 108, 194, 152, 82, 135, 125, 170, 83, 145, 53, 182,}, types.Key{32, 165, 59, 6, 170, 141, 121, 224, 241, 46, 75, 213, 187, 230, 93, 166, 153, 90, 219, 68, 45, 231, 178, 228, 122, 34, 43, 133, 29, 252, 241, 7}},
		},

		{
			{types.Key{115, 226, 61, 109, 227, 250, 2, 112, 90, 140, 86, 81, 26, 174, 241, 154, 47, 231, 137, 122, 120, 238, 234, 220, 236, 219, 50, 154, 39, 177, 122, 6,}, types.Key{72, 33, 192, 16, 86, 250, 1, 157, 76, 176, 231, 49, 89, 135, 103, 29, 2, 239, 105, 254, 160, 43, 248, 219, 81, 237, 251, 148, 180, 105, 39, 25}},
			{types.Key{176, 221, 118, 108, 183, 83, 221, 127, 242, 146, 30, 39, 176, 162, 233, 82, 138, 245, 147, 98, 13, 35, 137, 194, 96, 176, 209, 188, 10, 18, 117, 181,}, types.Key{204, 31, 27, 222, 33, 246, 87, 202, 235, 53, 102, 74, 75, 1, 158, 163, 237, 31, 95, 227, 249, 186, 221, 210, 249, 160, 110, 13, 253, 133, 8, 200}},
			{types.Key{17, 139, 231, 160, 130, 57, 92, 71, 14, 72, 142, 99, 102, 201, 49, 229, 180, 107, 134, 202, 168, 186, 101, 156, 79, 149, 167, 216, 72, 70, 183, 80,}, types.Key{229, 19, 71, 242, 30, 91, 215, 141, 230, 243, 74, 121, 146, 118, 100, 92, 150, 201, 171, 128, 219, 203, 190, 98, 94, 90, 4, 14, 186, 52, 18, 62}},
			{types.Key{237, 4, 226, 33, 232, 81, 213, 204, 76, 164, 151, 195, 158, 255, 247, 190, 22, 161, 225, 63, 195, 161, 231, 119, 26, 2, 98, 103, 156, 255, 247, 241,}, types.Key{173, 119, 64, 116, 17, 240, 128, 10, 116, 227, 192, 228, 112, 28, 249, 246, 74, 46, 76, 16, 209, 238, 133, 205, 235, 67, 153, 160, 113, 4, 146, 205}},
			{types.Key{55, 129, 166, 8, 80, 80, 115, 44, 60, 227, 229, 67, 218, 74, 179, 126, 70, 209, 162, 40, 194, 185, 53, 103, 64, 5, 71, 28, 7, 240, 155, 128,}, types.Key{123, 21, 184, 174, 209, 211, 220, 12, 205, 255, 213, 98, 136, 173, 179, 112, 253, 194, 9, 240, 49, 27, 21, 83, 169, 206, 246, 152, 217, 250, 43, 91}},
			{types.Key{117, 91, 215, 207, 65, 46, 252, 214, 87, 204, 246, 72, 204, 40, 157, 41, 106, 43, 162, 139, 44, 130, 166, 188, 104, 69, 134, 232, 68, 137, 211, 29,}, types.Key{191, 186, 89, 19, 22, 164, 126, 128, 151, 107, 153, 41, 234, 239, 164, 98, 175, 174, 221, 250, 2, 98, 90, 141, 126, 115, 220, 223, 224, 244, 135, 17}},
			{types.Key{199, 57, 137, 55, 39, 212, 164, 114, 133, 188, 245, 68, 191, 99, 112, 148, 151, 68, 155, 187, 110, 10, 157, 67, 41, 54, 142, 5, 49, 73, 145, 187,}, types.Key{167, 51, 65, 19, 173, 0, 186, 204, 186, 207, 12, 129, 149, 206, 20, 28, 100, 228, 17, 146, 30, 139, 248, 140, 76, 134, 115, 69, 4, 208, 242, 140}},
			{types.Key{160, 137, 245, 131, 0, 63, 85, 99, 108, 215, 95, 236, 134, 187, 174, 107, 87, 237, 59, 98, 66, 120, 113, 188, 73, 107, 232, 57, 3, 75, 52, 13,}, types.Key{111, 163, 46, 235, 215, 246, 45, 121, 175, 177, 170, 107, 230, 81, 189, 225, 212, 53, 116, 176, 92, 167, 192, 138, 216, 48, 203, 216, 156, 4, 84, 58}},
			{types.Key{137, 41, 126, 171, 48, 231, 34, 43, 86, 77, 218, 202, 122, 14, 62, 181, 22, 186, 157, 186, 149, 73, 131, 11, 187, 247, 53, 157, 31, 90, 232, 52,}, types.Key{129, 0, 224, 239, 84, 50, 73, 109, 189, 125, 161, 177, 77, 218, 240, 192, 108, 69, 188, 159, 88, 111, 81, 30, 54, 25, 105, 145, 100, 43, 247, 66}},
			{types.Key{146, 156, 138, 46, 195, 113, 106, 124, 30, 8, 10, 30, 47, 58, 177, 102, 23, 89, 241, 26, 63, 23, 209, 193, 33, 231, 116, 180, 86, 234, 169, 40,}, types.Key{245, 250, 115, 119, 36, 105, 215, 24, 160, 199, 210, 80, 245, 13, 1, 229, 54, 145, 218, 51, 22, 211, 62, 146, 140, 186, 212, 175, 164, 112, 229, 11}},
			{types.Key{141, 190, 151, 164, 205, 246, 60, 111, 167, 198, 230, 39, 16, 9, 44, 133, 146, 153, 240, 250, 14, 134, 238, 135, 102, 39, 30, 143, 66, 247, 81, 189,}, types.Key{160, 94, 18, 158, 51, 1, 68, 143, 186, 251, 196, 38, 179, 13, 54, 254, 214, 67, 99, 237, 75, 247, 94, 112, 110, 186, 24, 198, 93, 227, 236, 52}},
		},
	}

	pseudoOutsp = types.KeyV{sToKey("0x6e16f915e0c5c45f1569f201de7e28e6d667bf261cc8669a2b6cc228e312f3f7"), sToKey("0x229e3884291a9af85ec085131db25ed24c84042432332a6eb46ee5827ceb70a3")}

	bulletproof = types.Bulletproof{
		V:    types.KeyV{types.Key{121, 201, 148, 20, 165, 225, 8, 37, 186, 117, 239, 0, 3, 148, 76, 241, 86, 55, 38, 123, 182, 35, 115, 126, 76, 56, 186, 191, 23, 80, 177, 49}, types.Key{26, 147, 229, 167, 215, 242, 199, 47, 231, 16, 233, 227, 242, 178, 70, 89, 248, 248, 207, 138, 54, 17, 16, 176, 107, 247, 101, 177, 77, 58, 37, 148}},
		A:    sToKey("0x95647eb6cd4f4069772856993834342b7f690f59eed1f301f283e8b14405363f"),
		S:    sToKey("0x6b6165b145168b27067b5c0d818eed9a30cccde64992e54173e93dead100cf70"),
		T1:   sToKey("0x55c57232820320b96461891139e66ce1db32c619a60f2025bfc8c548541e8398"),
		T2:   sToKey("0x3bb482fa82c45c48b4eb97d6209e84a5f0cb162b89bba6faa6ec4e902584e90e"),
		Taux: sToKey("0xd66c860c3aefc2ab1b0cf313322301d7e8f3708fdc0c3357a96c90001cbd0304"),
		Mu:   sToKey("0x13bca9af218a7b338f545ce41207a569bc2782a287579c540bc5512fa1df160c"),
		L: types.KeyV{sToKey("0x0e8c4d421ba587014396c7b6a4d2d8f712e36b1e77ab6c7a7af392323730de47"), sToKey("0x076f2b793683ed05185a84f709f68e2d690cd23f173c34644a0d671e6b2c4f11"), sToKey("0xf8d7f6e2daab96fd7af323c6ee151726160e4909ca8daf41eacfca18284b7ae0"),
			sToKey("0xa739997085ca3edae4626e071f4fedce3932d588d3b50a3243ad477363c954f4"), sToKey("0x528c36c85da0f9d2bc9717d1dcc8c32d87dae876828d7a0efb613692165378fd"), sToKey("0x4a0aa4c991d3fb8d43ddc09ae38fe997a639021ad59809655e07822f800a294b"), sToKey("0x56432f98804f3cfcaf8e01b0754f6238edfb2a3acdccd9ce0d467c77002f6e24")},
		R: types.KeyV{sToKey("0xb6568f296b10a8d1013f2678318696d0e8ee789a8068d58132b3d23cf33ebfa1"), sToKey("0xaef046121dafb1ab0b413f9af87e2fa5d778ff1672ea5a4cd35816dcb3a785e5"), sToKey("0x9a6bd083d81492bd6eabba028453d9f135ecd89eed568986570c2e2d5ede9193"),
			sToKey("0x3a0e619c330e2e4115eb6d058bebf325488a41e0baf244079506f49a143c557e"), sToKey("0xc1e91c8453be0d937940612f0f12791bbf88b792997ba1992f5cccad2966f76a"), sToKey("0xcf789a52460ff1849068528f5703c06776003e7986b83377f1cedb990e18f1b9"), sToKey("0x931a749bbc369277d2fc5d19ff815c61c28868437ccc6e5e27e357227785086e")},
		Aa: sToKey("0x759d97bf74dd6d7764a33f81a8a0877a4d390937c531eefdedbb4ee7353d4905"),
		B:  sToKey("0x040a261a96e8505eec140dd99bc08be9feb8ab19eca1fa0a595671666f656b01"),
		T:  sToKey("0xfd8f60033729870acc31d32fb7a5592f869d2f697c80f0e1f8c3826f85559406"),
	}

	mgSigs = []types.MgSig{
		types.MgSig{
			Ss: types.KeyM{
				types.KeyV{sToKey("0xb0ebd47583a3ad53b2c5e81c9dd57db73c6edc4607832396f1f4ea90600d2f03"), sToKey("0x8d10c7b11cd3e490d32caeb0d6d8a7621f2cd7356ffc0378e6e706eeb4dcc80c")},
				types.KeyV{sToKey("0xdb28fd9f00488191bee2263995c447688b755c6ed311df1d04177ff064ecc90a"), sToKey("0x15b5176ee4624259aea0259ed3b21f7a62dc40d3afeef57321cd4caf007fbb07")},
				types.KeyV{sToKey("0x8b29beecfb269c92ed44d995f2340d21cf8b308f4a7359741c1702999c8cb00b"), sToKey("0x5f829dfb14dc6a5322b3b2225519f888d2623a8a71f7d02dfc5217917b5e5408")},
				types.KeyV{sToKey("0xcca9dfbfd690d6be02706e539d0e1adefbc861bf7b1ec326a0d301557600c108"), sToKey("0x67aa00b4dc7c5c819cfd888ebea6dd421f50ffab5275b46845166d59b7816707")},
				types.KeyV{sToKey("0x5e1259fa6088a8e69e1657c4614773a2323be77fae59d67d920989bece554305"), sToKey("0x696b91d3bd87dc054cc9fa9b67bb9cb5aad0ba68f7202f5a57141d730a43f707")},
				types.KeyV{sToKey("0xff6f07493bbaa2dd4a1b875e1ae171a87678c0c51f1dc56fb9b1dad531b88403"), sToKey("0xc847078f863fb9efd9ba55ba57821fe4e96566d2e8fc2ca28dc20d928ddf0907")},
				types.KeyV{sToKey("0xdf73ba3960b13e17918d8ba08bbfad13bd40b190ee9662ee5ba2d82c5ed5e70a"), sToKey("0x4ae985bbe2d05c141a9a9a68c99d6e44be113fbd158e3f18674bad5f101e810c")},
				types.KeyV{sToKey("0x2acd98b64ee1675eb122a6d19fab9ddb6e52096b53e8bb03ed4ec06ab55f3e0a"), sToKey("0x633fc00b33fc5bb4e91308065899531ee9366cdf0c392520234866fa7e5c1b00")},
				types.KeyV{sToKey("0x00fb5ecb25e9504e631b51d6b04f39c10145bd105e1a54f170cedc8e188ef20a"), sToKey("0x865245a3da0a41028da0e936bf2ae50b2a677363f00f6ffbcc2e79ecf284ad06")},
				types.KeyV{sToKey("0x5cac81d78877264e736fdbcb3b901d01f8733c30175eaecddc81d13b6238a101"), sToKey("0xed8faa8b51d710c35c3e6667e77376c990fae88bad4e7bce336abe6678b0a60f")},
				types.KeyV{sToKey("0x4a4b14da38fd62320fd2ce10c2139293850a4224d0266c8aa925fe1a6bc17806"), sToKey("0x1c9cbd1b1f0a0f1d93654f7aed566e77e636c3e51dde45aa9d5eb551ef609e09")},
			},
			Cc: sToKey("0x0a3186b128926e9e0eb312c53142b18f60ba8599e24684a1bc5246be7a86f808"),
			II: types.KeyV{types.Key{208, 41, 240, 108, 118, 93, 126, 9, 149, 201, 148, 139, 32, 180, 42, 139, 25, 44, 55, 222, 74, 158, 224, 75, 111, 223, 103, 61, 69, 227, 22, 166}},
		},
		types.MgSig{
			Ss: types.KeyM{
				types.KeyV{sToKey("0x13a89af5f38d8b17cd508c4cb2173697aeb23b704f078c442b721aa713c34809"), sToKey("0xf6369c5b253c09866cc370a21bd1c776de7ff3dc83313392cd3741be7387e703")},
				types.KeyV{sToKey("0x3a2a45491bbe76eaecbc640ef8048323fe056f362705288bc606dd7dbbd41e09"), sToKey("0x628217cb73aee63b97d624752340f011cde10f947532bb42aa3dd8ff639e970a")},
				types.KeyV{sToKey("0xa5e0d679b0aba29057ae7e3f9892311fd26a34f34862e72d55fda84db4ade10e"), sToKey("0xa02fab4f6d1e2055619bcc696b980966106576e16e8bea66588c81398fe5e302")},
				types.KeyV{sToKey("0xa570b7ba4137be56437cd7f3e1b79e2783e7e9f2f5313e7b20f705e14dbf4808"), sToKey("0x01d5d60b45bed35b78a1d24c6e0820d7849323afac56b88509f2a050af8a830e")},
				types.KeyV{sToKey("0x7a3d2c7dcfc51bdfd13c3ce7b9f0dd67f0a6dd87941fffca344f2f2e982cab06"), sToKey("0x6754d22db9a0859a79636bf679eb8822aa11a81e7510624395aa4eaa94fcac09")},
				types.KeyV{sToKey("0xed9a8b1fbcd385e72706f7d043650722398fa6e94f2994efae1bcd64d553350e"), sToKey("0x7a5e346b498b08f217b37314013b7355bad9cf2ff3a063e4b2962057b86aaf01")},
				types.KeyV{sToKey("0xb73f30ace6e20d1d6402dbffe05e053046e9ad1bd095877153b26290983e7f0c"), sToKey("0x2e55a7f0a83d6569e110444f84f2e79c56fbbc44159d2f1c04b291fa27341a07")},
				types.KeyV{sToKey("0x8e24c1701fb9e031c78591f717a1e3faf024b49c0c685cd3bcbddf4ca19f2805"), sToKey("0xf0c29cc3470f94dcda18b5e3cbe1bec892b633a54994a73fc5083aba5ed98603")},
				types.KeyV{sToKey("0x873a50b23bc2e3ed486cab83c612bc7bc943b82b0744f2849639f11789c63907"), sToKey("0x495a80fccb670da171ddc21befc169cebe45bd8cb2b28e0a99d2cb00af751a09")},
				types.KeyV{sToKey("0x7c453f220b9e4de474ef4e2443214707727b7e0825b96177261cb3bfdaa2d909"), sToKey("0x55246206cb3598c365410e1c462df41eeefec1b4aa347e0e5204f4e613816d0d")},
				types.KeyV{sToKey("0x61789c0472886b4a93e3cf53fda42f6a81a22c96a4f78104e8e4a7ef4baa8409"), sToKey("0x31d501ae8e06dfd216c204350df8f95bc69eb3e473c3553112a660bc26f9b407")},
			},
			Cc: sToKey("0x9cef8a2152d4109aeb4268fcb6cab98834c628c878271a086609e2716cea5c0a"),
			II: types.KeyV{types.Key{118, 14, 180, 13, 114, 66, 168, 22, 140, 145, 194, 226, 216, 96, 20, 92, 90, 131, 242, 255, 207, 166, 205, 204, 11, 151, 219, 5, 240, 174, 53, 136}},
		},
	}

	ecdhInfos = []types.EcdhTuple{
		types.EcdhTuple{
			Amount: types.Key{102, 35, 180, 221, 10, 77, 117, 238, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		types.EcdhTuple{
			Amount: types.Key{225, 24, 20, 225, 104, 225, 5, 138, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	outPk = types.CtkeyV{
		types.Ctkey{
			Dest: types.Key{225, 110, 138, 107, 13, 197, 103, 87, 227, 39, 246, 209, 177, 246, 188, 224, 71, 194, 37, 126, 41, 37, 24, 118, 173, 217, 155, 137, 86, 198, 191, 32},
			Mask: types.Key{22, 221, 175, 246, 221, 238, 44, 238, 240, 191, 7, 253, 181, 179, 131, 38, 197, 207, 195, 121, 166, 151, 159, 41, 85, 136, 212, 159, 151, 124, 55, 40},
		},
		types.Ctkey{
			Dest: types.Key{88, 87, 82, 176, 51, 147, 198, 174, 48, 114, 106, 81, 67, 28, 241, 217, 230, 218, 194, 6, 164, 186, 176, 83, 90, 67, 255, 61, 88, 155, 78, 167},
			Mask: types.Key{144, 125, 7, 195, 200, 209, 39, 203, 223, 50, 21, 168, 188, 207, 129, 166, 230, 67, 81, 165, 172, 133, 161, 231, 39, 102, 83, 194, 252, 115, 196, 33},
		},
	}
)

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
func hexdecode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("input is empty")
	}
	if has0xPrefix(input) {
		return hex.DecodeString(input[2:])

	} else {
		return hex.DecodeString(input)
	}
}

func sToKey(s string) (key types.Key) {
	x, err := hexdecode(s)
	if err != nil {
		panic(err.Error())
	}
	copy(key[:], x[:])
	return key
}

//for test
func DefualtTestRctsig() *types.RctSig {
	rctsig := &types.RctSig{}
	rctsig.Type = uint8(types.RCTTypeBulletproof2)
	rctsig.TxnFee = 140490000
	rctsig.Message = message
	rctsig.MixRing = mixring
	//rctsig.PseudoOuts = pseudoOutsp
	rctsig.EcdhInfo = ecdhInfos
	rctsig.OutPk = outPk
	rctsig.P.Bulletproofs = []types.Bulletproof{bulletproof}
	rctsig.P.MGs = mgSigs
	rctsig.P.PseudoOuts = pseudoOutsp
	return rctsig
}
func TestVerRctNonSemanticsSimple(t *testing.T) {
	if !VerRctNonSemanticsSimple(DefualtTestRctsig()) {
		t.Fatalf("TestVerRctNonSemanticsSimple fail");
	}
}
func TestVerRctSemanticsSimple(t *testing.T) {
	rctsig := DefualtTestRctsig()
	if !VerRctSemanticsSimple(rctsig) {
		t.Fatalf("TestVerRctSemanticsSimple fail");
	}
}
func TestVerRctSimple(t *testing.T) {
	rctsig := DefualtTestRctsig()
	if !VerRctSimpleTlv(rctsig) {
		t.Fatalf("TestVerRctSimple fail");
	}
}

func TestVerRct(t *testing.T) {
	//if !VerRct(DefualtTestRctsig()) {
	//	t.Fatalf("TestVerRct fail");
	//}
}
func TestProveRangeBulletproof(t *testing.T) {
	sk := types.KeyV{{59, 173, 27, 112, 149, 59, 198, 233, 243, 11, 237, 185, 50, 159, 46, 83, 39, 91, 225, 183, 137, 12, 66, 59, 133, 65, 187, 13, 248, 171, 85, 7}, {36, 173, 90, 206, 188, 160, 166, 80, 251, 173, 97, 206, 160, 31, 217, 251, 142, 12, 154, 59, 120, 35, 2, 253, 214, 248, 226, 249, 89, 238, 204, 14}}
	amounts := []types.Lk_amount{100000000, 35075991376858}
	_, c, mask, err := ProveRangeBulletproof(FromLkamountsToKeyv(amounts), sk)

	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(c) != 2 {
		t.Fatalf("C len expected 2,got %v", len(c))
	}
	if len(mask) != 2 {
		t.Fatalf("mask len expected 2,got %v", len(mask))
	}
	goodc := types.KeyV{sToKey("cae8b16ad94b692ed0b9e56cced83227b2eace578d1a9113a77e493bb51baf16"), sToKey("b4cf0efb1320aacd7e68a1090a03a0aadcf05f405404803926a91fae6abd8dd5")}
	goodmask := types.KeyV{sToKey("2ea2b81a53ed85cbb38e2a08a472156384254cf3f31580e72e3382a5eb957c00"), sToKey("7276523fb7304fbc16f558b5eb327b2db4c98b92c30ceb6baed5a65612cdc00b")}

	for i := 0; i < len(c); i++ {
		if !goodc[i].IsEqual(&c[i]) {
			fmt.Printf("C[%d]=%v\n", i, hex.EncodeToString(c[i][:]))
			t.Fatalf("C is invalid")
		}
	}
	for i := 0; i < len(mask); i++ {
		if !mask[i].IsEqual(&goodmask[i]) {
			fmt.Printf("mask[%d]=%v\n", i, hex.EncodeToString(mask[i][:]))
			t.Fatalf("mask is invalid")
		}
	}
	//fmt.Printf("Bulletproof A=%v\n",hex.EncodeToString(p.A[:]))
	//fmt.Printf("Bulletproof S=%v\n",hex.EncodeToString(p.S[:]))
	//fmt.Printf("Bulletproof T1=%v\n",hex.EncodeToString(p.T1[:]))
	//fmt.Printf("Bulletproof T2=%v\n",hex.EncodeToString(p.T2[:]))
}
func TestVerBulletproof(t *testing.T) {
	sk := types.KeyV{{59, 173, 27, 112, 149, 59, 198, 233, 243, 11, 237, 185, 50, 159, 46, 83, 39, 91, 225, 183, 137, 12, 66, 59, 133, 65, 187, 13, 248, 171, 85, 7}, {36, 173, 90, 206, 188, 160, 166, 80, 251, 173, 97, 206, 160, 31, 217, 251, 142, 12, 154, 59, 120, 35, 2, 253, 214, 248, 226, 249, 89, 238, 204, 14}}
	amounts := []types.Lk_amount{100000000, 35075991376858}
	bp, c, mask, err := ProveRangeBulletproof(FromLkamountsToKeyv(amounts), sk)

	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(c) != 2 {
		t.Fatalf("C len expected 2,got %v", len(c))
	}
	if len(mask) != 2 {
		t.Fatalf("mask len expected 2,got %v", len(mask))
	}
	goodc := types.KeyV{sToKey("cae8b16ad94b692ed0b9e56cced83227b2eace578d1a9113a77e493bb51baf16"), sToKey("b4cf0efb1320aacd7e68a1090a03a0aadcf05f405404803926a91fae6abd8dd5")}
	goodmask := types.KeyV{sToKey("2ea2b81a53ed85cbb38e2a08a472156384254cf3f31580e72e3382a5eb957c00"), sToKey("7276523fb7304fbc16f558b5eb327b2db4c98b92c30ceb6baed5a65612cdc00b")}

	for i := 0; i < len(c); i++ {
		if !goodc[i].IsEqual(&c[i]) {
			fmt.Printf("C[%d]=%v\n", i, hex.EncodeToString(c[i][:]))
			t.Fatalf("C is invalid")
		}
	}
	for i := 0; i < len(mask); i++ {
		if !mask[i].IsEqual(&goodmask[i]) {
			fmt.Printf("mask[%d]=%v\n", i, hex.EncodeToString(mask[i][:]))
			t.Fatalf("mask is invalid")
		}
	}

	ret, err := xcrypto.TlvVerBulletproof(bp)
	if err != nil {
		t.Fatalf("TlvVerBulletproof fail %v", err.Error())
	}
	if ret == false {
		t.Fatalf("TlvVerBulletproof fail")
	}
}

//TODO:
func TestProveRctMGSimple(t *testing.T) {
	//ProveRctMGSimple
}
func TestGetPreMlsagHashTlv(t *testing.T) {
	key, err := GetPreMlsagHash(DefualtTestRctsig())
	if err != nil {
		t.Fatalf(err.Error());
	}
	expectKey := sToKey("0x43ba268dfd93a4b8a491b95d63fec46e660424c180399af50e67e495948d4025")
	if !key.IsEqual(&expectKey) {
		t.Fatalf("unexpect key  got=%v,want=%v", hexEncode(key[:]), hexEncode(expectKey[:]))
	}
}
func TestGetPreMlsagHash(t *testing.T) {
	key, err := xcrypto.GetPreMlsagHash(DefualtTestRctsig())
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectKey := sToKey("0x43ba268dfd93a4b8a491b95d63fec46e660424c180399af50e67e495948d4025")
	if !key.IsEqual(&expectKey) {
		t.Fatalf("unexpect key  got=%v,want=%v", hexEncode(key[:]), hexEncode(expectKey[:]))
	}
}
func BenchmarkGetPreMlsagHashTlv(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetPreMlsagHash(sig)
		if err != nil {
			b.Fatalf("BenchmarkGetPreMlsagHashTlv fail");
		}
	}
}
func BenchmarkGetPreMlsagHash(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xcrypto.GetPreMlsagHash(sig)
	}
}
func BenchmarkVerRctNonSemanticsSimple(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerRctNonSemanticsSimple(sig) {
			b.Fatalf("TestVerRctNonSemanticsSimple fail");
		}
	}
}
func BenchmarkVerRctSemanticsSimple(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerRctSemanticsSimple(sig) {
			b.Fatalf("BenchmarkVerRctSemanticsSimple fail");
		}
	}
}
func BenchmarkVerRctSimple(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerRctSimple(sig) {
			b.Fatalf("BenchmarkVerRctSimple fail");
		}
	}
}

func BenchmarkVerRctSimpleTlv(b *testing.B) {
	sig := DefualtTestRctsig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerRctSimpleTlv(sig) {
			b.Fatalf("BenchmarkVerRctSimpleTlv fail");
		}
	}
}

func BenchmarkProveRangeBulletproof(b *testing.B) {
	sk := types.KeyV{{59, 173, 27, 112, 149, 59, 198, 233, 243, 11, 237, 185, 50, 159, 46, 83, 39, 91, 225, 183, 137, 12, 66, 59, 133, 65, 187, 13, 248, 171, 85, 7}, {36, 173, 90, 206, 188, 160, 166, 80, 251, 173, 97, 206, 160, 31, 217, 251, 142, 12, 154, 59, 120, 35, 2, 253, 214, 248, 226, 249, 89, 238, 204, 14}}
	amounts := []types.Lk_amount{100000000, 35075991376858}
	amountv := FromLkamountsToKeyv(amounts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProveRangeBulletproof(amountv, sk)
	}
}
func BenchmarkTlvVerBulletproof(b *testing.B) {
	sk := types.KeyV{{59, 173, 27, 112, 149, 59, 198, 233, 243, 11, 237, 185, 50, 159, 46, 83, 39, 91, 225, 183, 137, 12, 66, 59, 133, 65, 187, 13, 248, 171, 85, 7}, {36, 173, 90, 206, 188, 160, 166, 80, 251, 173, 97, 206, 160, 31, 217, 251, 142, 12, 154, 59, 120, 35, 2, 253, 214, 248, 226, 249, 89, 238, 204, 14}}
	amounts := []types.Lk_amount{100000000, 35075991376858}
	amountv := FromLkamountsToKeyv(amounts)
	pb, _, _, err := ProveRangeBulletproof(amountv, sk)
	if err != nil {
		b.Fatalf(err.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xcrypto.TlvVerBulletproof(pb)
	}
}
