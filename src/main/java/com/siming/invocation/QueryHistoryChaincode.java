package com.siming.invocation;

import com.siming.client.CaClient;
import com.siming.client.ChannelClient;
import com.siming.client.FabricClient;
import com.siming.config.Config;
import com.siming.user.UserContext;
import com.siming.util.Util;
import org.hyperledger.fabric.sdk.*;

import java.security.InvalidParameterException;
import java.util.Collection;
import java.util.logging.Level;
import java.util.logging.Logger;

public class QueryHistoryChaincode {

    public static void main(String[] args) {

        if(args.length != 1) {
            throw new InvalidParameterException("parameter must be exact 1.");
        }

        try {
            Util.cleanUp();
            String ca_URL = Config.CA_ORG1_URL;
            CaClient caClient = new CaClient(ca_URL, null);
            UserContext adminContext = new UserContext();
            adminContext.setName(Config.ADMIN);
            adminContext.setMspId(Config.ORG1_MSP);
            adminContext.setAffiliation(Config.ORG1);
            caClient.setCaContext(adminContext);
            adminContext = caClient.enrollAdminUser(Config.ADMIN,Config.ADMIN_PASSWORD);

            FabricClient fabricClient = new FabricClient(adminContext);
            ChannelClient channelClient = fabricClient.createChannelClient(Config.TRACE_CHANEL);
            Channel channel = channelClient.getChannel();
            Peer peer = fabricClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
            EventHub eventHub = fabricClient.getInstance().newEventHub("eventhub01","grpc://localhost:7053");
            Orderer orderer = fabricClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
            channel.addPeer(peer);
            channel.addEventHub(eventHub);
            channel.addOrderer(orderer);
            channel.initialize();

            Logger.getLogger(QueryChaincode.class.getName()).log(Level.INFO, "Query Product");
            Collection<ProposalResponse> responsesQuery = channelClient.queryByChainCode(Config.CHAINCODE_TRACE_NAME, "getProductHistory", args);
            for (ProposalResponse pres : responsesQuery) {
                String stringResponse = new String(pres.getChaincodeActionResponsePayload());
                Logger.getLogger(QueryChaincode.class.getName()).log(Level.INFO, stringResponse);
            }

        } catch (Exception e) {
            e.printStackTrace();
        }




    }
}
