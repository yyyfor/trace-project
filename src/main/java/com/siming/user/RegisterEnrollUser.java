package com.siming.user;

import com.siming.client.CaClient;
import com.siming.config.Config;
import com.siming.util.Util;

public class RegisterEnrollUser {

    public static void main(String[] args) {

        try {
            Util.cleanUp();
            String ca_URL = Config.CA_ORG1_URL;
            CaClient caClient = new CaClient(ca_URL, null);
            UserContext adminUserContext = new UserContext();
            adminUserContext.setName(Config.ADMIN);
            adminUserContext.setMspId(Config.ORG1_MSP);
            adminUserContext.setAffiliation(Config.ORG1);
            caClient.setCaContext(adminUserContext);
            adminUserContext = caClient.enrollAdminUser(Config.ADMIN, Config.ADMIN_PASSWORD);

            UserContext userContext = new UserContext();
            String name = "user" + System.currentTimeMillis();
            userContext.setName(name);
            userContext.setAffiliation(Config.ORG1);
            userContext.setMspId(Config.ORG1_MSP);
            String secret = caClient.registerUser(name, Config.ORG1);
            userContext = caClient.enrollUser(userContext, secret);

        } catch (Exception e) {
            e.printStackTrace();
        }



    }
}
