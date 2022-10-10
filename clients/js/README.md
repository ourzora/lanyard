# lanyard

lanyard is a javascript client for [lanyard.org](https://lanyard.org) – a decentralized way to create and consume allowlists.

http is handled by [isomorphic-fetch](https://www.npmjs.com/package/isomorphic-fetch); this package targets browsers and backend servers.

## functionality

### creating an allowlist

```js
import lanyard from 'lanyard'

const resp = await lanyard.createTree({
  unhashedLeaves: [
    '0xfb843f8c4992efdb6b42349c35f025ca55742d33',
    '0x7e5507281f62c0f8d666beaea212751cd88994b8',
    '0xd8da6bf26964af9d7eed9e03e53415d37aa96045',
  ],
  // leafTypeDescriptor: ["address"] // optional, used for abi encoded types
  // packedEncoding: boolean // optional, default false
})

console.log(resp.merkleRoot)
// 0x8aeeaf632a31342dfccb7dd4f1654ec602c263b33769062bd6ed59d1644d2af6
```

### getting a tree by root

```js
const tree = await lanyard.getTree(
  '0x8aeeaf632a31342dfccb7dd4f1654ec602c263b33769062bd6ed59d1644d2af6',
)

console.log(tree)
// {
//   "unhashedLeaves": [
//     "0xfb843f8c4992efdb6b42349c35f025ca55742d33",
//     "0x7e5507281f62c0f8d666beaea212751cd88994b8",
//     "0xd8da6bf26964af9d7eed9e03e53415d37aa96045"
//   ],
//   "leafCount": 3,
//   "leafTypeDescriptor": null,
//   "packedEncoding": false
// }
```

### getting a proof for item

```js
const proof = await lanyard.getProof({
  merkleRoot:
    '0x8aeeaf632a31342dfccb7dd4f1654ec602c263b33769062bd6ed59d1644d2af6',
  unhashedLeaf: '0xfb843f8c4992efdb6b42349c35f025ca55742d33',
})

console.log(proof)
// {
//   "unhashedLeaf": "0xfb843f8c4992efdb6b42349c35f025ca55742d33",
//   "proof": [
//     "0xdb740d4f5f900a98f8513824cbcb164917f4e0b948914b750613b76063b70565",
//     "0x06e120c2c3547c60ee47f712d32e5acf38b35d1cc62e23b055a69bb88284c281"
//   ]
// }
```

### getting roots for a proof

```js
const roots = await lanyard.getRoots([
  '0xdb740d4f5f900a98f8513824cbcb164917f4e0b948914b750613b76063b70565',
  '0x06e120c2c3547c60ee47f712d32e5acf38b35d1cc62e23b055a69bb88284c281',
])

console.log(roots)

// {
//   "roots": [
//     "0x8aeeaf632a31342dfccb7dd4f1654ec602c263b33769062bd6ed59d1644d2af6"
//   ]
// }
```
