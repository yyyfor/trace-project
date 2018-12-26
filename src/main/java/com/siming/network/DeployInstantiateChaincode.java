package com.siming.network;

import com.siming.client.ChannelClient;
import com.siming.client.FabricClient;
import com.siming.config.Config;
import com.siming.user.UserContext;
import com.siming.util.Util;
import org.hyperledger.fabric.protos.peer.Chaincode;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.CryptoException;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric.sdk.TransactionRequest.Type;


import java.io.File;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

public class DeployInstantiateChaincode {

    public static void main(String[] args) {
        try {
            CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();
            UserContext org1Admin = new UserContext();
            File pkFolder1 = new File(Config.ORG1_USR_ADMIN_PK);
            File [] pkFiles1 = pkFolder1.listFiles();
            File certFolder = new File(Config.ORG1_USR_ADMIN_CERT);
            File [] certFiles1 = certFolder.listFiles();
            Enrollment enrollOrg1Admin = Util.getEnrollment(Config.ORG1_USR_ADMIN_PK, pkFiles1[0].getName(),
                    Config.ORG1_USR_ADMIN_CERT, certFiles1[0].getName());
            org1Admin.setEnrollment(enrollOrg1Admin);
            org1Admin.setMspId(Config.ORG1_MSP);
            org1Admin.setName(Config.ADMIN);

            UserContext org2Admin = new UserContext();
            File pkFolder2 = new File(Config.ORG2_USR_ADMIN_PK);
            File[] pkFiles2 = pkFolder2.listFiles();
            File certFolder2 = new File(Config.ORG2_USR_ADMIN_CERT);
            File[] certFiles2 = certFolder2.listFiles();
            Enrollment enrollOrg2Admin = Util.getEnrollment(Config.ORG2_USR_ADMIN_PK, pkFiles2[0].getName(),
                    Config.ORG2_USR_ADMIN_CERT, certFiles2[0].getName());
            org2Admin.setEnrollment(enrollOrg2Admin);
            org2Admin.setMspId(Config.ORG2_MSP);
            org2Admin.setName(Config.ADMIN);

            UserContext org3Admin = new UserContext();
            File pkFolder3 = new File(Config.ORG3_USR_ADMIN_PK);
            File[] pkFiles3 = pkFolder3.listFiles();
            File certFolder3 = new File(Config.ORG3_USR_ADMIN_CERT);
            File[] certFiles3 = certFolder3.listFiles();
            Enrollment enrollOrg3Admin = Util.getEnrollment(Config.ORG3_USR_ADMIN_PK, pkFiles3[0].getName(),
                    Config.ORG3_USR_ADMIN_CERT, certFiles3[0].getName());
            org3Admin.setEnrollment(enrollOrg3Admin);
            org3Admin.setMspId(Config.ORG3_MSP);
            org3Admin.setName(Config.ADMIN);

            FabricClient fabricClient = new FabricClient(org1Admin);
            Channel trace_channel = fabricClient.getInstance().newChannel(Config.CHANNEL_NAME);
            Orderer orderer = fabricClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
            Peer peer0_org1 = fabricClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
            Peer peer1_org1 = fabricClient.getInstance().newPeer(Config.ORG1_PEER_1, Config.ORG1_PEER_1_URL);
            Peer peer0_org2 = fabricClient.getInstance().newPeer(Config.ORG2_PEER_0, Config.ORG2_PEER_0_URL);
            Peer peer1_org2 = fabricClient.getInstance().newPeer(Config.ORG2_PEER_1, Config.ORG2_PEER_1_URL);
            Peer peer0_org3 = fabricClient.getInstance().newPeer(Config.ORG3_PEER_0, Config.ORG3_PEER_0_URL);
            Peer peer1_org3 = fabricClient.getInstance().newPeer(Config.ORG3_PEER_1, Config.ORG3_PEER_1_URL);
            trace_channel.addOrderer(orderer);
            trace_channel.addPeer(peer0_org1);
            trace_channel.addPeer(peer1_org1);
            trace_channel.addPeer(peer0_org2);
            trace_channel.addPeer(peer1_org2);
            trace_channel.addPeer(peer0_org3);
            trace_channel.addPeer(peer1_org3);
            trace_channel.initialize();

            List<Peer> org1Peers = new ArrayList<Peer>();
            List<Peer> org2Peers = new ArrayList<Peer>();
            List<Peer> org3Peers = new ArrayList<Peer>();
            org1Peers.add(peer0_org1);
            org1Peers.add(peer1_org1);
            org2Peers.add(peer0_org2);
            org2Peers.add(peer1_org2);
            org3Peers.add(peer0_org3);
            org3Peers.add(peer1_org3);

            Collection<ProposalResponse> response = fabricClient.deployChainCode(Config.CHAINCODE_1_NAME,Config.CHAINCODE_1_PATH,
                    Config.CHAINCODE_ROOT_DIR, Chaincode.ChaincodeSpec.Type.GOLANG.toString(), Config.CHAINCODE_1_VERSION, org1Peers);

            for (ProposalResponse res : response) {
                Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
                        Config.CHAINCODE_1_NAME + "- Chain code deployment " + res.getStatus());
            }

            fabricClient.getInstance().setUserContext(org2Admin);

            response = fabricClient.deployChainCode(Config.CHAINCODE_1_NAME,
                    Config.CHAINCODE_1_PATH, Config.CHAINCODE_ROOT_DIR, TransactionRequest.Type.GO_LANG.toString(),
                    Config.CHAINCODE_1_VERSION, org2Peers);

            for (ProposalResponse res : response) {
                Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
                        Config.CHAINCODE_1_NAME + "- Chain code deployment " + res.getStatus());
            }

            fabricClient.getInstance().setUserContext(org3Admin);

            response = fabricClient.deployChainCode(Config.CHAINCODE_1_NAME,
                    Config.CHAINCODE_1_PATH, Config.CHAINCODE_ROOT_DIR, TransactionRequest.Type.GO_LANG.toString(),
                    Config.CHAINCODE_1_VERSION, org3Peers);

            for (ProposalResponse res : response) {
                Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
                        Config.CHAINCODE_1_NAME + "- Chain code deployment " + res.getStatus());
            }

            ChannelClient channelClient = new ChannelClient(trace_channel.getName(), trace_channel, fabricClient);
            String [] arguments = {""};
            response = channelClient.instantiateChainCode(Config.CHAINCODE_1_NAME, Config.CHAINCODE_1_VERSION,
                    Config.CHAINCODE_1_PATH, Type.GO_LANG.toString(), "init", arguments, null);
            for (ProposalResponse res : response) {
                Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
                        Config.CHAINCODE_1_NAME + "- Chain code instantiation " + res.getStatus());
            }

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
