package com.siming.client;

import com.siming.user.UserContext;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.CryptoException;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

import java.io.File;
import java.lang.reflect.InvocationTargetException;
import java.util.Collection;
import java.util.logging.Level;
import java.util.logging.Logger;

public class FabricClient {

    private HFClient instance;

    public HFClient getInstance() {
        return instance;
    }

    public FabricClient(UserContext context) throws CryptoException, InvalidArgumentException, ClassNotFoundException, NoSuchMethodException, InvocationTargetException, InstantiationException, IllegalAccessException {
        CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();
        instance = HFClient.createNewInstance();
        instance.setCryptoSuite(cryptoSuite);
        instance.setUserContext(context);
    }


    public ChannelClient createChannelClient(String name) throws InvalidArgumentException {
        Channel channel = instance.newChannel(name);
        ChannelClient client = new ChannelClient(name, channel, this);
        return client;
    }

    public Collection<ProposalResponse> deployChainCode(String chaincodeName, String chaincodePath, String codepath,
                                                        String language, String version, Collection<Peer> peers) throws InvalidArgumentException, ProposalException {
        InstallProposalRequest request = instance.newInstallProposalRequest();
        ChaincodeID.Builder chaincodeIDBuilder = ChaincodeID.newBuilder().setName(chaincodeName).setVersion(version).setPath(chaincodePath);
        ChaincodeID chaincodeID= chaincodeIDBuilder.build();
        Logger.getLogger(FabricClient.class.getName()).log(Level.INFO,
                "Deploying chaincode " + chaincodeName + " using Fabric client " + instance.getUserContext().getMspId()
                        + " " + instance.getUserContext().getName());
        request.setChaincodeID(chaincodeID);
        request.setUserContext(instance.getUserContext());
        request.setChaincodeSourceLocation(new File(codepath));
        request.setChaincodeVersion(version);
        Collection<ProposalResponse> responses = instance.sendInstallProposal(request, peers);
        return responses;
    }


}
