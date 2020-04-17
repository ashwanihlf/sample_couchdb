# Hyperledger Fabric network setup with With CouchDB to explore chaincode queries  

Pre Requisite - Hyperledger Binaries and HLF Pre-Requisites software are installed

# Following are the steps to run the setup
1. create a working folder, change directory to working folder
2. git clone https://github.com/ashwanihlf/sample_couchdb.git
3. sudo chmod -R 755 sample_couchdb/
4. cd sample_couchdb
5. mkdir config  
	<remove config and crypto-config if they are existing before creation of config folder (Optional)>
	5a. sudo rm -rf config
	5b  sudo rm -rf crypto-config
6. export COMPOSE_PROJECT_NAME=net
7. sudo ./generate.sh
8. sudo ./start.sh
9. docker exec -it cli /bin/bash
10. peer chaincode invoke -C mychannel -n samplecc -c '{"Args":["initLedger"]}'
11. peer chaincode invoke -C mychannel -n samplecc -c '{"Args":["getAllCars"]}'

>> returns result: status:200 payload:"[
{\"Key\":\"CAR0\", \"Record\":{\"color\":\"blue\",\"docType\":\"\",\"model\":\"bmw\",\"owner\":\"Ashwani\"}}
{\"Key\":\"CAR1\", \"Record\":{\"color\":\"red\",\"docType\":\"\",\"model\":\"santro\",\"owner\":\"Raja\"}}
{\"Key\":\"CAR10\", \"Record\":{\"color\":\"silver\",\"docType\":\"\",\"model\":\"thar\",\"owner\":\"Yoyo\"}}
{\"Key\":\"CAR11\", \"Record\":{\"color\":\"gold\",\"docType\":\"\",\"model\":\"xuv\",\"owner\":\"Mini\"}}
{\"Key\":\"CAR2\", \"Record\":{\"color\":\"white\",\"docType\":\"\",\"model\":\"wagonR\",\"owner\":\"Naman\"}}
{\"Key\":\"CAR3\", \"Record\":{\"color\":\"black\",\"docType\":\"\",\"model\":\"fortuner\",\"owner\":\"Gurneet\"}}
{\"Key\":\"CAR4\", \"Record\":{\"color\":\"orange\",\"docType\":\"\",\"model\":\"pajero\",\"owner\":\"John\"}}
{\"Key\":\"CAR5\", \"Record\":{\"color\":\"grey\",\"docType\":\"\",\"model\":\"RangeRover\",\"owner\":\"Bob\"}}
{\"Key\":\"CAR6\", \"Record\":{\"color\":\"voilet\",\"docType\":\"\",\"model\":\"bentley\",\"owner\":\"Alice\"}}
{\"Key\":\"CAR7\", \"Record\":{\"color\":\"ruby\",\"docType\":\"\",\"model\":\"merc\",\"owner\":\"Mayank\"}}]"

12. peer chaincode invoke -C mychannel -n samplecc -c '{"Args":["getCarsByRange","CAR0","CAR5"]}'


>> returns result: status:200 payload:"[
{\"Key\":\"CAR0\", \"Record\":{\"color\":\"blue\",\"docType\":\"\",\"model\":\"bmw\",\"owner\":\"Ashwani\"}}
{\"Key\":\"CAR1\", \"Record\":{\"color\":\"red\",\"docType\":\"\",\"model\":\"santro\",\"owner\":\"Raja\"}}
{\"Key\":\"CAR10\", \"Record\":{\"color\":\"silver\",\"docType\":\"\",\"model\":\"thar\",\"owner\":\"Yoyo\"}}
{\"Key\":\"CAR11\", \"Record\":{\"color\":\"gold\",\"docType\":\"\",\"model\":\"xuv\",\"owner\":\"Mini\"}}
{\"Key\":\"CAR2\", \"Record\":{\"color\":\"white\",\"docType\":\"\",\"model\":\"wagonR\",\"owner\":\"Naman\"}}
{\"Key\":\"CAR3\", \"Record\":{\"color\":\"black\",\"docType\":\"\",\"model\":\"fortuner\",\"owner\":\"Gurneet\"}}
{\"Key\":\"CAR4\", \"Record\":{\"color\":\"orange\",\"docType\":\"\",\"model\":\"pajero\",\"owner\":\"John\"}}]"
