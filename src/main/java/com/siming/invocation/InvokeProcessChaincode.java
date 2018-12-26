package com.siming.invocation;

import com.siming.client.CaClient;
import com.siming.client.ChannelClient;
import com.siming.client.FabricClient;
import com.siming.config.Config;
import com.siming.user.UserContext;
import com.siming.util.Util;
import org.hyperledger.fabric.sdk.*;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;

import static java.nio.charset.StandardCharsets.UTF_8;

public class InvokeProcessChaincode {

    private static final byte[] EXPECTED_EVENT_DATA = "!".getBytes(UTF_8);
    private static final String EXPECTED_EVENT_NAME = "event";

    public static void main(String[] args) {

        try {
            Util.cleanUp();
            String caUrl = Config.CA_ORG1_URL;
            CaClient caClient = new CaClient(caUrl, null);
            UserContext adminUserContext = new UserContext();
            adminUserContext.setName(Config.ADMIN);
            adminUserContext.setMspId(Config.ORG1_MSP);
            adminUserContext.setAffiliation(Config.ORG1);
            caClient.setCaContext(adminUserContext);
            adminUserContext = caClient.enrollAdminUser(Config.ADMIN, Config.ADMIN_PASSWORD);
            FabricClient fabricClient = new FabricClient(adminUserContext);

            ChannelClient channelClient = fabricClient.createChannelClient(Config.TRACE_CHANEL);
            Channel channel = channelClient.getChannel();
            Peer peer = fabricClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
            EventHub eventHub = fabricClient.getInstance().newEventHub("eventhub01", "grpc://localhost:7053");
            Orderer orderer = fabricClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
            channel.addOrderer(orderer);
            channel.addPeer(peer);
            channel.addEventHub(eventHub);
            channel.initialize();

            TransactionProposalRequest request = fabricClient.getInstance().newTransactionProposalRequest();
            ChaincodeID ccid = ChaincodeID.newBuilder().setName(Config.CHAINCODE_TRACE_NAME).setVersion(Config.CHAINCODE_TRACE_VERSION).build();
            request.setChaincodeID(ccid);
            request.setFcn("getProduct");
            String[] arguments = { "0"};
            request.setArgs(arguments);
//            request.setFcn("getIdHistory");
//            String[] arguments = { "5"};
//            request.setFcn("getInfo");
//            String[] arguments = { "5"};
//            request.setArgs(arguments);
            request.setProposalWaitTime(1000);

            Map<String, byte[]> tm2 = new HashMap<>();
            tm2.put("HyperLedgerFabric", "TransactionProposalRequest:JavaSDK".getBytes(UTF_8));
            tm2.put("method", "TransactionProposalRequest".getBytes(UTF_8));
            tm2.put("result", ":)".getBytes(UTF_8));
            tm2.put(EXPECTED_EVENT_NAME, EXPECTED_EVENT_DATA);
            request.setTransientMap(tm2);
            Collection<ProposalResponse> responses = channelClient.sendTransactionProposal(request);
            for (ProposalResponse res: responses) {
                ChaincodeResponse.Status status = res.getStatus();
                Logger.getLogger(InvokeProcessChaincode.class.getName()).log(Level.INFO,"Invoked putVal on "+Config.CHAINCODE_1_NAME + ". Status - " + status);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

    }
}
