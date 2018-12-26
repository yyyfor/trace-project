#!/bin/bash
#
# Exit on first error, print all commands.
set -e

#Start from here
echo -e "\nStopping the previous network (if any)"
docker-compose -f docker-compose.yml down

# If need to re-generate the artifacts, uncomment the following lines and run
#
#cryptogen generate --config=./crypto-config.yaml
#mkdir config
#configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./config/genesis.block
#configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./config/channel.tx -channelID trace-channel 
#configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID trace-channel -asOrg Org1MSP
#configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID trace-channel -asOrg Org2MSP
#configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/Org3MSPanchors.tx -channelID trace-channel -asOrg Org3MSP
# Create and Start the Docker containers for the network
echo -e "\nSetting up the Hyperledger Fabric 1.1 network"
docker-compose -f docker-compose.yml up -d
sleep 15
echo -e "\nNetwork setup completed!!\n"

