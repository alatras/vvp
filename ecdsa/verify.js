const Accounts = require('web3-eth-accounts')
const accounts = new Accounts('ws://localhost:8546')

function main() {
  const message = process.argv[2]
  const address = process.argv[3]
  const signature = process.argv[4]
  console.log(accounts.recover(message, signature))
  if (accounts.recover(message, signature) == address)
    console.log("\nVERIFICATION PASSED")
  else
    console.log("\nVERIFICATION FAILED")
}

if (require.main === module) main()