const { runTypeChain, glob } = require('typechain')
const { execSync } = require('child_process')
const { unlinkSync } = require('fs')

const getPackagePath = packageName => {
  const path = require.resolve(`${packageName}/package.json`)
  return path.substr(0, path.indexOf('package.json'))
}

async function main() {
  const cwd = process.cwd()

  const nitroPath = getPackagePath('@arbitrum/nitro-contracts')
  const tokenBridgePath = getPackagePath('@arbitrum/token-bridge-contracts')

  console.log('Compiling paths.')

  const npmExec = process.env['npm_execpath']
  if (!npmExec || npmExec === '')
    throw new Error(
      'No support for npm_execpath env variable in package manager'
    )

  // TODO: use `HARDHAT_ARTIFACT_PATH` to write files to arbitrum sdk instead of the packages themselves.
  // this is currently broken since hardhat throws a weird error:
  // `Error HH702: Invalid artifact path [...] its correct case-sensitive path is...`
  // https://yarnpkg.com/advanced/rulebook#packages-should-never-write-inside-their-own-folder-outside-of-postinstall
  // instead of writing in postinstall in each of those packages, we should target a local folder in sdk's postinstall

  // copy the hardhat config to nitro-contracts
  execSync(
    `cp ${cwd}/hardhat.nitro-contracts-abigen.ts ${nitroPath}/hardhat-abigen.ts`
  )
  // copy the hardhat config to token-bridge-contracts
  execSync(
    `cp ${cwd}/hardhat.token-bridge-contracts-abigen.ts ${tokenBridgePath}/hardhat-abigen.ts`
  )

  console.log('building @arbitrum/nitro-contracts')
  execSync(`${npmExec} run build --config hardhat-abigen.ts`, {
    cwd: nitroPath,
  })

  console.log('building @arbitrum/token-bridge-contracts')
  execSync(`${npmExec} run build --config hardhat-abigen.ts`, {
    cwd: tokenBridgePath,
  })

  console.log('Done compiling')

  const nitroFiles = glob(cwd, [
    `${tokenBridgePath}/build/contracts/!(build-info)/**/+([a-zA-Z0-9_]).json`,
    `${nitroPath}/build/contracts/!(build-info)/**/+([a-zA-Z0-9_]).json`,
  ])

  // TODO: generate files into different subfolders (ie `/nitro/*`) to avoid overwrite of contracts with the same name
  await runTypeChain({
    cwd,
    filesToProcess: nitroFiles,
    allFiles: nitroFiles,
    outDir: './src/lib/abi/',
    target: 'ethers-v5',
  })

  const classicFiles = glob(cwd, [
    // we have a hardcoded abi for the old outbox
    `./src/lib/dataEntities/Outbox.json`,
  ])

  await runTypeChain({
    cwd,
    filesToProcess: classicFiles,
    allFiles: classicFiles,
    outDir: './src/lib/abi/classic',
    target: 'ethers-v5',
  })

  // we delete the index file since it doesn't play well with tree shaking
  unlinkSync(`${cwd}/src/lib/abi/index.ts`)
  unlinkSync(`${cwd}/src/lib/abi/classic/index.ts`)

  console.log('Typechain generated')
}

main()
  .then(() => console.log('Done.'))
  .catch(console.error)
