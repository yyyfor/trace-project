package com.siming.network;

import com.siming.client.FabricClient;
import com.siming.config.Config;
import com.siming.user.UserContext;
import com.siming.util.Util;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

import java.io.File;
import java.util.Collection;
import java.util.Iterator;
import java.util.logging.Level;
import java.util.logging.Logger;

public class CreateChannel {

    public static void main(String[] args) {
        try {
            CryptoSuite.Factory.getCryptoSuite();
            Util.cleanUp();
            UserContext org1Admin = new UserContext();
            File pkFolder1 = new File(Config.ORG1_USR_ADMIN_PK);
            File [] pkFiles1 = pkFolder1.listFiles();
            File certFolder1 = new File(Config.ORG1_USR_ADMIN_CERT);
            File [] certFiles1 = certFolder1.listFiles();
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
            File [] pkFiles3 = pkFolder3.listFiles();
            File certFolder3 = new File(Config.ORG3_USR_ADMIN_CERT);
            File [] certFiles3 = certFolder3.listFiles();
            Enrollment enrollOrg3Admin = Util.getEnrollment(Config.ORG3_USR_ADMIN_PK, pkFiles3[0].getName(),
                    Config.ORG3_USR_ADMIN_CERT, certFiles3[0].getName());
            org3Admin.setEnrollment(enrollOrg3Admin);
            org3Admin.setMspId(Config.ORG3_MSP);
            org3Admin.setName(Config.ADMIN);

            FabricClient fabricClient = new FabricClient(org1Admin);

            Orderer orderer = fabricClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
            ChannelConfiguration channelConfiguration = new ChannelConfiguration(new File(Config.CHANNEL_CONFIG_PATH));
            byte [] channelConfigurationSignatures = fabricClient.getInstance().getChannelConfigurationSignature(channelConfiguration, org1Admin);
            Channel trace_channel = fabricClient.getInstance().newChannel(Config.CHANNEL_NAME, orderer, channelConfiguration, channelConfigurationSignatures);
            Peer peer0_org1 = fabricClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
            Peer peer1_org1 = fabricClient.getInstance().newPeer(Config.ORG1_PEER_1, Config.ORG1_PEER_1_URL);
            Peer peer0_org2 = fabricClient.getInstance().newPeer(Config.ORG2_PEER_0, Config.ORG2_PEER_0_URL);
            Peer peer1_org2 = fabricClient.getInstance().newPeer(Config.ORG2_PEER_1, Config.ORG2_PEER_1_URL);
            Peer peer0_org3 = fabricClient.getInstance().newPeer(Config.ORG3_PEER_0, Config.ORG3_PEER_0_URL);
            Peer peer1_org3 = fabricClient.getInstance().newPeer(Config.ORG3_PEER_1, Config.ORG3_PEER_1_URL);
            trace_channel.joinPeer(peer0_org1);
            trace_channel.joinPeer(peer1_org1);
            trace_channel.addOrderer(orderer);
//            trace_channel.initialize();
            fabricClient.getInstance().setUserContext(org2Admin);
            trace_channel = fabricClient.getInstance().getChannel(Config.CHANNEL_NAME);
            trace_channel.joinPeer(peer0_org2);
            trace_channel.joinPeer(peer1_org2);
            fabricClient.getInstance().setUserContext(org3Admin);
            trace_channel = fabricClient.getInstance().getChannel(Config.CHANNEL_NAME);
            trace_channel.joinPeer(peer0_org3);
            trace_channel.joinPeer(peer1_org3);
            Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO, "Channel created "+trace_channel.getName());
            Collection peers = trace_channel.getPeers();
            Iterator peerIter = peers.iterator();
            while(peerIter.hasNext()) {
                Peer pr = (Peer) peerIter.next();
                Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO,pr.getName()+ " at " + pr.getUrl());
            }

        } catch (Exception e) {
            e.printStackTrace();
        }

    }

}
