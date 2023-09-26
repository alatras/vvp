const crypto = require('@polkadot/util-crypto')

async function main() {
  await crypto.cryptoWaitReady()
  const message = process.argv[2]
  const publicKey = process.argv[3]
  const signature = process.argv[4]
  if ((await crypto.signatureVerify(message, signature, publicKey)).isValid)
    console.log("\nVERIFICATION PASSED")
  else
    console.log("\nVERIFICATION FAILED")
}

if (require.main === module) main()