package main

import (
	"encoding/hex"
	"math"
	"testing"
)

type testData struct {
	input          string
	entropyPerByte float64
	chiSquare      float64
	mean           float64
}

/* sample data generated using the ent program (https://packages.ubuntu.com/jammy/ent) */
var data = []testData{
	{input: "7baef0770d0d28e296b92bf28458cf43265a339eed0c9b5f17fad5a7b1f6aef4312998524c95a9edded79da322f7ace3b703404229c8a58edea3f9f774f3513c", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 141.078125},
	{input: "772f568472bf57a94ebfbf2305f855f9fb07ab02dba23a4ebec2c04969f235f0db5b038d71f24dfd028d182699d791ff12732e7d724e98172ca3e607a0e4a1dc", entropyPerByte: 5.663910, chiSquare: 288.000000, mean: 129.062500},
	{input: "2dfed1a03f56cb36f5e0142eaafe5e19ae0e908f64edbd77b024fc4109daf0a88ac6a6bf1d621df787edfec02c8760a5b5f86a1c43bb10fed1811300e2f4487e", entropyPerByte: 5.750000, chiSquare: 272.000000, mean: 139.093750},
	{input: "e85a0dd7164db59d6fd0fe6931e7c290ce01e6daec1e0d1f5877006232bc0970b3d9b551c4e795424d65a649b2c8b22ba4a538f8afd26167af9029b850d09424", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 132.218750},
	{input: "f40ad4e71c00a8e1decdd5da26c29860bf3d9bd8b1bb6defe803e7f60a2bc9c3628786b58e7061314ba7f0cf1e6bcd385e9508bde56fd7b4e9d845b779d945f1", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 148.609375},
	{input: "d2c301949df45b0db3fb6af667d5526f87076e176b148b92ddbc8aa8649ecec271d1084710e1478e281e18f990e618b6feb90b9e2efbf4d31fa69f77637c708c", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 132.781250},
	{input: "e135caa75f556161fddb02a143b9006ed3fc79fc539d3062f04b0f10a1a5599f5964806faddcff181b15f0e8729c6c0106a53dd11160c9cd5f20b2b332ed7f05", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 125.343750},
	{input: "d2686fe6ff0ceba4f77e47b35d6e962ac8d19cd9279f7a0980b86ccbeeef497c411e8a8dbb6a78054438c35326bb8d7ca0731e40dfb5887a224476e53371055d", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 127.703125},
	{input: "b40a6f8a24e02a0718febfe679adbe2e16048473eebb67f418e11c83052f3a453c0bcab7e4dfec8d3d0991dc1c71cdfdd78f04ff39b25a1e962db226fdfe55b8", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 126.343750},
	{input: "b048748a022595ff125dd552edc2352cb2392115ccd28f9eca624193c1da104bad00f6507742f1b7c1f4345cd7e9206184ae37d4221b5ded313e41dae6a1e463", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 128.734375},
	{input: "b03d34ef6067b4a372deea5d4ff61e924c87a221c205328e7e30194221e331e79aa08dc1721f4a685e060167227615c9a1169aea730079143709ccaf179965b4", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 111.125000},
	{input: "69d956cb4efb63ba3e943eb234aad50774391fc56a9771ef2a0481ff1c51a6fcbf8a163614cec69513e129c95c2c969cd5dbf1041e46cab7784d55dc20c69de9", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 129.281250},
	{input: "d3957f62d4c5c70eadc42bc9be16fcc10277d88f0d18dd5807dc96bc92db3a39e2d4424c382e9cc51cc56551f1ee8951ea4ce952676a68ad85c5d11e4ddaf394", entropyPerByte: 5.750000, chiSquare: 272.000000, mean: 138.328125},
	{input: "6fc0c4debe3d4b14b4ed03191dbf22c7497d8addc27359e4f93ef8c14eb19eea839b41b1c58dc585cf484a81c5125b6551d5f87b040f488d28b9a66bc8d24fdb", entropyPerByte: 5.800705, chiSquare: 248.000000, mean: 136.484375},
	{input: "b77531eaeb3cbc0c478f7afdc9b897037b29f564d3e8d025057252408528721d518cebb2c39388c56d4e0a63d43db1fead77f70435fa91b4d40021b7f4284c41", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 129.953125},
	{input: "a40463ec017a54e25a2a848e0729f540a3581f56047b461b48823507a7ae6785a4113652b963ba717498a98b9f5f61730dc82d0f6dce9dc4970963c1f1cbbd00", entropyPerByte: 5.831955, chiSquare: 240.000000, mean: 110.906250},
	{input: "2db142d2f3cc3d7e2e4a8627ad21ec1884dbc3007837c9ec9effc81a722e7ff28c8b027fee9dd3ec12dfe24c79c03304a7f29ac972c4f4daa1c283747a0c0e9e", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 137.843750},
	{input: "1175155a325c5307a4cdc9555c58c681eb9a5ee15872c1712d3ffc1e4f6330b5602b2b147555c289957e2991860bb98167b1cfafea99de5ea9ab58af6bf35f6a", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 122.281250},
	{input: "356d74221a1f7534097f4c257d13f2d55ef09cba3121e393f7ef78053243852cc7d8473ee232cab55bf6f2a3a902b96c1b1fe9d75d3a50e0626f3cdad5376f62", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 121.265625},
	{input: "991c8f297d97db259e0219d80da6f17f24d0d00b85a59f1cbc927698cc2334c510940b3b29153234f8e914d9347ba06d81127193309300bc447e059f7d35ce7a", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 110.921875},
	{input: "16e0fca335a2be2a93e3472cc7923796be7ef534250d84f123568c9488af613d9732b5844d86f8f770bb7876bbf642ee7d23b09a2a66873a7746ee932ee513cb", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 132.078125},
	{input: "d442b8924cd301954aced73b8b565a07fa584f9486a06c9e6b53e4f03c79667c48d2a6ece095831ab6694813ad0e7cd92d7717db89a73519e1b324318c591d1b", entropyPerByte: 5.906250, chiSquare: 216.000000, mean: 121.765625},
	{input: "a86b96cd3ab184677cc42506c43904cecd2498f1c1da5fee3732eec5c64e2adce1f65195426099da2868dff8a484421c916520bb791479020e7ea3a27088163b", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 127.250000},
	{input: "0ad83716009ffdaa488dd9bddcd2640130c87568e7e87e1fea44766ea19aab5e2a1aefe7a6654dc6ce240bd11ba23e14de104fdc09a565b758e1387d1e181787", entropyPerByte: 5.906250, chiSquare: 216.000000, mean: 121.062500},
	{input: "eab5912e0b24219176eb75ddf87fe25d9b2aa8c4fe0d80521c7ea8f0245e65ef32860221f535b8eaeb3699568e3955a858ccc718b2e7037ac3b2243a37df2aba", entropyPerByte: 5.663910, chiSquare: 288.000000, mean: 128.484375},
	{input: "eaaf849b61c97385119428ffb1f94b64b101dff13e823cf15201feb84de676c87b52ea5ccdfaeb7fb0c17c1baefeaa055df201834b5af9761dff9fb58f6c721f", entropyPerByte: 5.644455, chiSquare: 288.000000, mean: 141.781250},
	{input: "5b9f394c235be662d8e06030f04666496a0535b569d8a559bb73e607e1e7ddca66aee9c52d41cf1756051ed86890e0136dec3993e3fedccaed293c156bc202ff", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 132.609375},
	{input: "572583d00d2924a77c2a0a7d72a2859ae28e80a3b1267e31b7335e36a5c9d36df1d7f8e0798c996ef179ca0fddb77489bad42ca5d958bee6b3ccfdf05e1d3651", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 138.359375},
	{input: "efac0d1fa99b7866f110f5740a0689c569fb00a124fef90863e1f50b1c4c850768290716fbe423a9887b0c6282b1a5ca82b2f2ef6e430b970af697c0023d55e1", entropyPerByte: 5.687500, chiSquare: 272.000000, mean: 122.765625},
	{input: "c416413e5ab214b00f1700e27b7ce75064bc03b98e240892a4f75ac2f6ae4cc1286915d5cf091df3bd43e8ac7b4cbcb23dd38b7451cd50486b8c5565f7543f37", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 121.281250},
	{input: "d160401e5edd2273bc2a6839a390365b2082eb611b4c7c084029a0af004e35eb12b3e4e298ed2c70c11899d7124a81898c8ef28f78073e6082fac2a3168999d2", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 118.609375},
	{input: "72bf6489900b11cfd5fba6d52307c2030a4a246da9131d0cdda8e446edf5002954f06e4cd32f5f840bd59a7bfe8f1bc02b5af4d6f0beb59a2d2aa87ea95a26bf", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 125.125000},
	{input: "4602e536050d694bbd8e3b2f917e2bd1cd90a700fcbb05c396f59ec194642635d317d2b0982ed52c7fddc65e7974bff0d2daf63429c7849b59a8981cb408e83b", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 129.765625},
	{input: "4fa9de6ae6c38e10bb1b52f8d0752b37d1780ecd516aa5252e80456b75f7cd9f53d8ea32a853a3da51a0d96c9809135bb4a65f7264d186cc404017ffa1c69ce0", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 133.437500},
	{input: "88097b0479794c67ee176335e4895b2a3965f54688d858b1d3bfdc240fee468c82661fa29ad054ffd8162481616a9221cfd3cd491613bb836d2916e85b086494", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 117.687500},
	{input: "f5f82476852b2d42e00f365d1cf34577236f39bf927c361cd6bd127dd042bc9af62db272b3692ee29c5d5f4e9daadf5119b92be80b21ae2e04117053f90ba17e", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 115.484375},
	{input: "1212254198fdb3aeb35b173647b9d55b0b18db2bc30f1353efd0ea77c6b8d5e0a970e2b6564e6f2631e009af9f89f5e478d07bb1df22d3f8df584aacd18680a1", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 137.375000},
	{input: "dacdac3a197df078410ec13961d475ebf4445cee380f36ea1d6cde4cb82d771c546cd7afb14638630c80132d47b304393266014edf6c37f1d50bcc41562c383b", entropyPerByte: 5.757660, chiSquare: 264.000000, mean: 110.312500},
	{input: "cc527a4b5f75e8e1ad281f042db5d6b296488b404726ec0a366fccf7f11dda805ba6e572a1f6b88be248a20542ba2a93700cd802c0032548354589eedbefe010", entropyPerByte: 5.863205, chiSquare: 232.000000, mean: 126.671875},
	{input: "a9d1ac83ddbd160cb2cfb4e275c10f2abd52629698877e4901993cff224a642a9cb307ba178a6826257de2e840e98818c52f449ea601509bcc6fed962332d089", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 123.734375},
	{input: "0dc1e56e62669ac2079a67e6fe23b68044baaf0535f462a1b72fa004cd97ec0cc45f9ea374644f10f038d62f3664b765b54e42e3003ca6e1344679674eb26d96", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 124.312500},
	{input: "f612ee4d5913e5b1251530a89154248bd5267650705d301ac0f7793fbdce65abfb9afce1de29a402977065e188be7dfa0c1391be9538aed89759c6c16abf548e", entropyPerByte: 5.687500, chiSquare: 272.000000, mean: 132.890625},
	{input: "97105283e8607cdc51a8c9f3f692703cd85376299c66d5f245665e7e8aa3279e9f7f9ad18d2dc4ca37e819983474b859a48cc7d34680a41182cb32481f5a6625", entropyPerByte: 5.863205, chiSquare: 232.000000, mean: 129.921875},
	{input: "837039e733bf48379dcbfcf359978acc37dd24362b4487be6b3b00a9bf0e079ed705569decab7f4421dde9fdfc55ffc7cf008b71337399e47b116db6858ebda5", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 132.906250},
	{input: "e8d1b625f897ee7c11adca954ff8c4a3e30359bae1fba9163e060a43104f372a14e0da964d24e650b74c22481699fca5edb0a57b3b46dedbf63df9dcab6fb822", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 135.875000},
	{input: "6028712b903a8680c5a902735472a733122d57464fc3a267aabc89ed926bb65aff24ac6f24b54c5c5703046f9db732e6ceec3ecce52dcb8f257fad0be27c4b61", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 119.578125},
	{input: "2622359f45f853b9e5ad5901fa993307ba9a7d66f86d883733149a00ce1458337546bc7dddc3396070fcee167bd429148263b9ab2fc727879fdc6f301f6a11c9", entropyPerByte: 5.695160, chiSquare: 280.000000, mean: 117.406250},
	{input: "d34477c715a52636192db0f30686b70d7283108f396081e62701e191af995e99053ca9603d5133965656fd7416fe1c818856ea406557bcd9dd61725e2abb860b", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 113.453125},
	{input: "84ba7e2a074898987ecbcdb61794e8f74f6724046917ce49f9f7e8334f953d1db82feb8f84076353ca1cd680009a5156f848cd72bde8802a2f9e2d10b0df351c", entropyPerByte: 5.519455, chiSquare: 320.000000, mean: 121.250000},
	{input: "1b850618b86db747a1652d85477277b02257017ecea3f9671cff97f043f377b2effe3979ad1b704ea1beadb3460d4aafecc34d0ecc78a6f18e19062694edb833", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 126.359375},
	{input: "9f35ad6f3b80509e7863807c7907bef6ae63c9d68d76021a995f9d9f96417a331ab829b20a18a7c703aeffac5a5e228edaaeb2f7f2db4f253c6da74ff0372528", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 122.765625},
	{input: "c963d31b9bf0b6c067634ff72b719e679da741d60b301432ebb379ad521867e6ef4b1d1deca7cf28fd04bd7a17826708936270489e5dbc54d120eae7adbc04b6", entropyPerByte: 5.625000, chiSquare: 304.000000, mean: 127.968750},
	{input: "0f777fbbaa227bf7612cdc7eb307fcad49e8f9aa73a21b81b620f8666c865172c7cb7d5136591344375bf3e3e5288086dbd1e915bc2af40ba05d78dac4cdc649", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 135.828125},
	{input: "8649152f01eb9ee2f43ca8b4516f89c72950c6ada26187a3944b6d480f16890e71a77dad221bd5852f3f0e1e08382bf43eb1317a1bf60c6f925f4892587d9e93", entropyPerByte: 5.656250, chiSquare: 280.000000, mean: 109.203125},
	{input: "d56da712b0b617e646f3aca5868d80e081c7e8345922844786d1c532a25fd4b0eefd26e7a0fcfccc190f60304881052e9d2a258848204d0526ea084faa3e255d", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 123.609375},
	{input: "d721843fe9311c41984a94e7121fee7e8440f5579f10f727ca8895fefb8ae69975a062671da93042fad60d2c9f03375099bcb80fc3323dfeba3bd338aad94264", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 126.734375},
	{input: "3fbf698a8e3baaa7b5ea315b7e8f65efdc7d1534ef80cfe1a4a9fcd4e2aa99a95d4bf1915f5c4c0c64db066adc93b415e8314406a6466234d2f4bbb7ce1bd34d", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 138.203125},
	{input: "936ab438d0b6eae7b0e082c6c51f550d9e707741a6483861ceaca2d95ff2abf7b22a0108e335b5ffb9300a5419096f0868b7ab1dc45b6ee31751aa231318530f", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 120.953125},
	{input: "07a69522d717dda6966e820ecb8a0b62c6ccbf0b23ef0cdcf52614e5bc246b9102a18923722794702771b97121583d0c6d344cfb895b91cfc53198dcd26b5ca6", entropyPerByte: 5.644455, chiSquare: 288.000000, mean: 118.843750},
	{input: "5c29dee7de6aa6bd9c2900a49bf6d9ee1e1b9e832cd1d878c7b9b376f787a2c7bee71c7040a5e94b43be846b39229c966bd171c5d9930afa1a16394e04735dc6", entropyPerByte: 5.687500, chiSquare: 272.000000, mean: 135.609375},
	{input: "2cd4a7dfa6989e07fe6e1afa20da1f2e7fd3782d84a06926f280f0a143d42cd07578235d307239640d77156d20f51286c3a31988f06a5f10b1dff4d585886347", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 124.734375},
	{input: "b0f962962bc439db1ade641b876e0d28cf555f29d89e4ae0e8c8b84447626444672560c1d944d7a0b281e7f4d8e475c76d98edbf483218dbb7acbdc33d7e3d6c", entropyPerByte: 5.769455, chiSquare: 256.000000, mean: 136.906250},
	{input: "634b588fb18721927978c659457438be76eb4d30078a312972b9472ef7a19aecbe87c64b38e0b8e84f4ba793563ffe7710656c4d4ea967d05181023446764129", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 116.937500},
	{input: "e650800cf6e25a5762db7dda7f524d95afb1efa97a3b1fa580ff4b7210723ab8dfd4ce056e4bbb13eddd2e918256ffc5a9861edf4b4950bd69c167da79d78a78", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 137.734375},
	{input: "068923c0ba52720dabd985789bc9b04cdd6b1e73ad19e14568156607ba7e920821248236457bd855813ab500fffed48618b793179cfa5dc5a78255d184d689af", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 124.531250},
	{input: "58c4b895dddf8040b6443e0b0559234d65a9f87d1901f90c2375e8353e19a7573b9f44dbd1f905a84ba43876327a109067fccf74cf0e1250f373eca22da1ca3a", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 118.593750},
	{input: "085bd496ff6ed22e1e0c1217e46689f86d6b3eeda5b40b4d78f9e6ee93659da0deb891e969d774e25f5768d014c1f39d46e4dc19b7ba49c656775dd8c292e1c0", entropyPerByte: 5.937500, chiSquare: 208.000000, mean: 144.328125},
	{input: "6b0e821a55f5884e02433d81ce343eab5b07f4da9787f68845d95ae33cbc0c1d6af90d6b77fd133592e77a8440e18e195b7bf3896b96a5f466288ec21a2b319b", entropyPerByte: 5.769455, chiSquare: 256.000000, mean: 120.437500},
	{input: "2bcb3fb4bfa7de31759f6dd0ea967edb30272aac5047505f8c778c1751645bdb3510ac96fac2f0b620ce6f3cbf95307b8fe61a2767d8762d77d4e9e05719e58c", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 130.406250},
	{input: "191733244d936675f450ccae0c35cdcf661019a21e35f782cce6d7243cc46216a966046022c9b2d1c7d0052928b572844d7f28ff3fdbad2dac5acfec664a7382", entropyPerByte: 5.625000, chiSquare: 304.000000, mean: 119.078125},
	{input: "d9d1e86785f553607086b4cbf2f9a2d95df171d011e07af46329a6c34c5583a01667140589418b89c0367b36d292658a7325d490a5525b392c3459dd6a95d868", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 133.781250},
	{input: "77b6b3d653e9e4abdfc9d5a50b145b772ec3e02a29ff0e563d8f779be976690ef0acf3f0b0660ca7f01c633842f526e7995df99372a4b2948bf303f41a91bf18", entropyPerByte: 5.757660, chiSquare: 264.000000, mean: 139.359375},
	{input: "d66759908030e2614c7bb3eb4a373cfc9fa0bc48a47ed4ff016c7153181798468b748164ea6e521233b9c93a6a76800b8295c4db82caeabb0761c9b341f69dff", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 132.078125},
	{input: "b8b15bd067b75c990db3e01f41fcd9e7e97ab23e778b2deedd164aac74512f96d439aa6816183fdc6eeb1dd7f5d536251606022d5600661791361c095e384594", entropyPerByte: 5.863205, chiSquare: 232.000000, mean: 114.390625},
	{input: "31e5c78ec4261e1c7ba25404d83eff15df5df4ff74592eaa34b62d4edbeb04c3deeea613e1f6b76ceeead7e968f63ad3d3bb3db2b2eeba21b227c048cec3cd09", entropyPerByte: 5.695160, chiSquare: 280.000000, mean: 146.843750},
	{input: "ef57d943742c358323adb44774a6f2182457cec69054fb1bf4e2df5fd3485e30063329d381447b8ab85074065c24e49014f98eb2c5fee19a287bc1e9d4c8fa92", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 135.562500},
	{input: "4a215393916dfd01db3bdb63cbdf93e74d486fa1123c6a1d29989539d7bf66b9720c0e48f197809a04204969e8f0bf088d53b35be45ef657b792f9d46d185f7e", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 125.546875},
	{input: "e286aeb14f2550fe347d376c78c7ec60f38e4ad5b4d72187426e347c0d2666a2bb4652ec2dab6f814f1e2aaf9eb0d38cc2d9cf878ad5734a5abc422f4dd12514", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 126.921875},
	{input: "5db531aaf298e0927b0d25f20b2283371d268663f24111773f44ff8b33a4d909f2a94c94691bb9cdaad5ad8703b74bc69d15488ac8dc5e65136382d187f876aa", entropyPerByte: 5.738205, chiSquare: 280.000000, mean: 125.781250},
	{input: "2ca728874b1d7dadd4c99369ae8179a8d51f5f91e3c42710b70961a6163fd2bc624756742918908db2a0990935fbda0b0e7ee41b7477f94bc68053c3efd9699a", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 124.562500},
	{input: "e6bd4aee54195d5cffaaea11d06777313eb00edf00f4e58329272818b42ead38ef4426fc2ba231465efb8db75aca3e999b3312391b5524efeaf6fa9760559e9c", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 124.984375},
	{input: "d8e97c40dd3ebb8dd7f3fe2a8d13b178b078c146ff3da04e99b9b400b5c9c2926d9f6e179ba060c1866d927a8bcfc76a2f296497b556375a5d151d90fbc1307d", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 134.906250},
	{input: "fa2871569f878ae14cf058d1fdec1dd1fba4fe699fd5943d301174c06dfd9d0eb3ee6f168f1824af267dfc730cb305afbed522241ff51038714addd8d30c5343", entropyPerByte: 5.718750, chiSquare: 264.000000, mean: 131.156250},
	{input: "f7fbce74a28dde95d56a083d22baa1cc403084846eb6e98575530d957e09d21e97b387067c5be953c7d40b40cb1ff88960fa250de491edd9081c3f6698cb854f", entropyPerByte: 5.718750, chiSquare: 264.000000, mean: 130.250000},
	{input: "5d1592b3871e3fd525b21a26f590c094ac54ddaeddf077dfeb01b6a3861d2770d4aac88b77dbd4eeab7a5ac8efee4aa1b93fa06356356dc65df857ee35603bb6", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 141.406250},
	{input: "3b7b2999686ef588f371a78b3cf6c322ee68f5f65177efa27cbce559b7777b6acda46887584108244563edb16a9dbd67ee0f907b37765d859db7f43b7ecc1b4d", entropyPerByte: 5.601410, chiSquare: 304.000000, mean: 135.953125},
	{input: "db99c30c65edb667b362873ab39527cb5b3d783cd01c6dd0250868cd8eeb7e11c3a5ddd4cbaf0c3ea9dcebdbd40594bd51fbaf6bc42e8e306e1fa55a43137798", entropyPerByte: 5.656250, chiSquare: 280.000000, mean: 131.906250},
	{input: "f66172e9f3ec1c70080a19a2e5407bf6e0b53697f2b065486db883a9eaea2c3cccfa04100e4d2b078ef0bcb49cfb4163bb2437b9b519f617c432785e52a94488", entropyPerByte: 5.800705, chiSquare: 248.000000, mean: 130.171875},
	{input: "3c13da45c413b873a2c08ec891c2a858edd7edd8fb28fb1520711f105efd4601e3eac2dc427e375f524e647c3ef859f9d4d5043ad9a82a6f1aef0635e746234b", entropyPerByte: 5.812500, chiSquare: 240.000000, mean: 129.031250},
	{input: "e3d31688871ab3f908071b2f2ce532299ee4b5d8ce5dc31076be753df6d94c80e47c41b8ce03831dc30e271b103c5acdf69564d08453ca79f45e3e278c1ebf7d", entropyPerByte: 5.781250, chiSquare: 248.000000, mean: 122.906250},
	{input: "9eb9656c311315c4e8cc144a1d7bf2162e5acfebbd40648a3ced728aeaab0b04195a6a462009444b00fea479f2199af43cfc9e78cace4a8df2e5568b391a11a1", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 119.843750},
	{input: "0ea794e41cc19a95c753538c42ac8b8ab9acaca5ca803547ce9f59cc478ad35d3d47e27bd896394c2194a67a02244deb9a93f4a7f59bf76f3bf4e655e4105595", entropyPerByte: 5.570160, chiSquare: 312.000000, mean: 137.265625},
	{input: "209236478f04230d4e251fcbdb5b8fe3cc023f3a882bf35f609d0571d60527ab30e6b2621d17e40c2f5eaea98ef30ef40fef12bca13ccb81a1deeacd93e145ba", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 118.843750},
	{input: "7349de31f23ce12a312b6416b84fc5f77b6b9aa6ec869d241e14c978aecda2c4a2fd458040f435f09b2826d48fa94302194adb097ccdc3e764aa5fa9ee65388c", entropyPerByte: 5.843750, chiSquare: 232.000000, mean: 130.625000},
	{input: "a682b5397a463b98c519d67d77d1bb74e4a1ab26fa43af41ec607820557facaf8c3a05e6d7f4c1894db85f020331a251284f18b70dfe99d053c9a2e9461b8878", entropyPerByte: 5.875000, chiSquare: 224.000000, mean: 127.937500},
	{input: "40e35fc401ed77ff7077cfca9cd0489229eb94786288b7487e25729409b2cf16081d91850a8bb4a512b31ee94f619f091a62071b09526084cbe31c33c74ea530", entropyPerByte: 5.706955, chiSquare: 272.000000, mean: 115.296875},
	{input: "409b327b49be83bb951f1914f068067583f548f2dd061eb1f7e76bc2bae02614447935ea5133c180c6144fba5ed28b3e6b1dae4707c8f9b577b58ae95b939191", entropyPerByte: 5.738205, chiSquare: 264.000000, mean: 127.656250},
	{input: "e93d5b3d01557bb6cd9001a386dae63ccc0c05f770613bb1dfd9908d7ed3f6f2055b1dc5e46c05e9cfaa2aadd52e1c895584bb3cc1b6bae34c7f5cef576dc226", entropyPerByte: 5.675705, chiSquare: 280.000000, mean: 134.078125},
	{input: "72662c6aba8a9dd1d1cce9defd31acd7c911fe0fb95f5c2a895892738458ba74a8eca3ecbdaa4efd2be5c607cc2b38168a9035566e2e4fcda04fcc876da377c6", entropyPerByte: 5.613205, chiSquare: 296.000000, mean: 138.203125},
	{input: "7a9f150b0e1904610c7a008b659411c35389e99d80ddaf4b872eef3a6672baa1cc079ec2e0e3db2ec4dfc73294c246475981f39f6c659207a1c570e19b3414a6", entropyPerByte: 5.750000, chiSquare: 256.000000, mean: 124.015625},
}

func almostEqual(a float64, b float64) bool {
	epsilon := 0.000001
	return math.Abs(a-b) < epsilon
}

func TestTracker_EntropyQA(t *testing.T) {
	var tracker Tracker
	for _, d := range data {
		input, _ := hex.DecodeString(d.input)
		epb := tracker.entropyPerByte(input)
		if !almostEqual(d.entropyPerByte, epb) {
			t.Errorf("entropy per byte calculation incorrect for %s, except: %f, get: %f", d.input, d.entropyPerByte, epb)
		}
		amd := tracker.arithmeticMeanDeviation(input)
		if !almostEqual(math.Abs(d.mean-127.5), amd) {
			t.Errorf("arithmetic mean deviation calculation incorrect for %s, except: %f, get: %f", d.input, math.Abs(d.mean-127.5), amd)
		}
		cs := tracker.chiSquare(input)
		if !almostEqual(d.chiSquare, cs) {
			t.Errorf("chi-square calculation incorrect for %s, except: %f, get: %f", d.input, d.chiSquare, cs)
		}
	}
}
