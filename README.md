All instructions below are composed by referring to https://hyperledger-fabric.readthedocs.io/en/latest/peer-chaincode-devmode.html

I just tried to simplify them.

---

__Fork this repo and then go to your GitHub account and clone the forked repo__

**Note: 

- Exec all commands in the root folder (in the outermost folder of the project or cloned repo).
- Do you want to use your chaincode?
    - Yes: Replace "rajan-31" in "cmd/main.go" with your GitHub username
    - No: Go ahead!

---

__Main Instructions__

1. Install binaries for orderer, peer, and configtxgen (there are many possible ways to achieve this. So, I didn't mention any instructions)

2. Only if you want to use your chaincode
    - Put chaincode in "my_chaincode" folder (put "main" function in cmd/main.go and other logic in a go file, in my case "fabcar.go")

3. Install go packages

    ```sh
    go mod init github.com/rajan-31/hyperledger-fabric_chaincode-devmode-contractapi
    go mod tidy
    ```

4. Set environment variables (if anywhere in the process like no config file found, then check if this env variable is set with `echo $FABRIC_CFG_PATH`)
    
    ```sh
    export FABRIC_CFG_PATH=$(pwd)/sampleconfig
    ```

5. Default location Fabric use to store blocks and other data (in 2nd command, use your username)
    
    ```sh
    sudo mkdir /var/hyperledger
    sudo chown rajan /var/hyperledger
    ```

6. Generate the genesis block for the ordering service
    
    ```sh
    configtxgen -profile SampleDevModeSolo -channelID syschannel -outputBlock genesisblock -configPath $FABRIC_CFG_PATH -outputBlock "$(pwd)/sampleconfig/genesisblock"
    ```

7. Start orderer

    ```sh
    export FABRIC_CFG_PATH=$(pwd)/sampleconfig

    ORDERER_GENERAL_GENESISPROFILE=SampleDevModeSolo orderer
    ```

8. Start the peer in DevMode (in a new terminal)

    ```sh
    export FABRIC_CFG_PATH=$(pwd)/sampleconfig

    FABRIC_LOGGING_SPEC=chaincode=debug CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052 peer node start --peer-chaincodedev=true
    ```

9. Create a channel and join (in a new terminal)

    ```sh
    export FABRIC_CFG_PATH=$(pwd)/sampleconfig
    
    configtxgen -channelID ch1 -outputCreateChannelTx ch1.tx -profile SampleSingleMSPChannel -configPath $FABRIC_CFG_PATH

    peer channel create -o 127.0.0.1:7050 -c ch1 -f ch1.tx

    peer channel join -b ch1.block
    ```

10. Build the chaincode

    ```sh
    # it will build the chaincode and put it in build/bin
    go build -o ./build/bin/smartContract ./my_chaincode/cmd
    ```

11. Start the chaincode in the root folder

    ```sh
    CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=false CORE_CHAINCODE_ID_NAME=mycc:1.0 ./build/bin/smartContract -peer.address 127.0.0.1:7052
    ```

12. Approve and commit the chaincode definition

    ```sh
    peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --package-id mycc:1.0

    peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')"

    peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --peerAddresses 127.0.0.1:7051
    ```

13. Query or Invoke chaincode

    ```sh
    export CORE_PEER_ADDRESS=127.0.0.1:7051

    peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Function":"InitLedger","Args":[]}' --isInit
    
    peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Function":"QueryAllCars","Args":[]}'

    peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Function":"QueryCar","Args":["CAR0"]}'

    peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Function":"CreateCar","Args":["CAR10","Mercedes-Benz", "EQS", "Black", "Rajan"]}'
    ```
\
\
_**Note: See "after-tutorial" branch to compare generated files, such as genesisblock, ch1.block, ch1.tx, go.mod, go.sum, smartContract_

---

__After following the above steps, to update your chaincode, you need to do is:__

1. Make any changes you want to chaincode (in my case to fabcar.go)

2. stop chaincode with "ctrl+c" (chaincode process that you started in step 11)

3. Repeat step 10 (to rebuild updated chaincode)

4. start chaincode again (follow step 11)

5. test chaincode (step 13)

*When you are done, stop orderer, peer, and chaincode.
\
\
\
__When you will again start development, you need to do is__

1. Start orderer (step 7)

2. Start peer in devmode (step 8)

3. Start Chaincode (step 11)