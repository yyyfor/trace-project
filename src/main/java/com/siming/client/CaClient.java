package com.siming.client;

import com.siming.user.UserContext;
import com.siming.util.Util;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.exception.CryptoException;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.hyperledger.fabric_ca.sdk.exception.EnrollmentException;

import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.net.MalformedURLException;
import java.util.Properties;
import java.util.logging.Level;
import java.util.logging.Logger;

public class CaClient {

    private String caURL;
    private Properties properties;

    private HFCAClient client;
    private UserContext caContext;

    public UserContext getCaContext() {
        return caContext;
    }

    public void setCaContext(UserContext caContext) {
        this.caContext = caContext;
    }

    public CaClient(String caURL, Properties properties) throws IllegalAccessException, InvocationTargetException, InvalidArgumentException, InstantiationException, ClassNotFoundException, NoSuchMethodException, MalformedURLException, CryptoException {
        this.caURL = caURL;
        this.properties = properties;
        init();
    }

    public void init() throws MalformedURLException, IllegalAccessException, InvocationTargetException, InvalidArgumentException, InstantiationException, NoSuchMethodException, CryptoException, ClassNotFoundException {
        CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();
        client = HFCAClient.createNewInstance(caURL, properties);
        client.setCryptoSuite(cryptoSuite);
    }

    public HFCAClient getClient() {
        return client;
    }

    public UserContext enrollAdminUser(String username, String password) throws IOException, ClassNotFoundException, EnrollmentException, org.hyperledger.fabric_ca.sdk.exception.InvalidArgumentException {
        UserContext userContext = Util.readUserContext(caContext.getAffiliation(), username);
        if (userContext != null) {
            Logger.getLogger(CaClient.class.getName()).log(Level.WARNING, "CA -" + caURL + " admin is already enrolled.");
            return userContext;
        }

        Enrollment adminEnrollment = client.enroll(username, password);
        caContext.setEnrollment(adminEnrollment);
        Logger.getLogger(CaClient.class.getName()).log(Level.INFO, "CA -" + caURL + " Enrolled Admin.");
        Util.writeUserContext(caContext);
        return caContext;
    }

    public String registerUser(String username, String organization) throws Exception {
        UserContext userContext = Util.readUserContext(caContext.getAffiliation(), username);
        if (userContext != null) {
            Logger.getLogger(CaClient.class.getName()).log(Level.WARNING, "CA -" + caURL +" User " + username+ " is already registered.");
            return null;
        }
        RegistrationRequest registrationRequest = new RegistrationRequest(username, organization);
        String enrollmentSecret = client.register(registrationRequest, caContext);
        Logger.getLogger(CaClient.class.getName()).log(Level.INFO, "CA -" + caURL + " Registered User - " + username);
        return enrollmentSecret;
    }

    public UserContext enrollUser(UserContext user, String secret) throws IOException, ClassNotFoundException, EnrollmentException, org.hyperledger.fabric_ca.sdk.exception.InvalidArgumentException {
        UserContext userContext = Util.readUserContext(caContext.getAffiliation(), user.getName());
        if (userContext != null) {
            Logger.getLogger(CaClient.class.getName()).log(Level.WARNING, "CA -" + caURL + " User " + user.getName()+" is already enrolled");
            return userContext;
        }
        Enrollment enrollment = client.enroll(user.getName(), secret);
        user.setEnrollment(enrollment);
        Util.writeUserContext(user);
        Logger.getLogger(CaClient.class.getName()).log(Level.INFO, "CA -" + caURL +" Enrolled User - " + user.getName());
        return user;
    }
}
