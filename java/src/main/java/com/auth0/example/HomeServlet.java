package com.auth0.example;

import com.auth0.SessionUtils;
import com.auth0.jwt.JWT;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

public class HomeServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse res)
            throws ServletException, IOException {

        String idToken = (String) SessionUtils.get(req, "idToken");

        String prettyJson = null;

        if (idToken != null) {
            DecodedJWT jwt = JWT.decode(idToken);

            Map<String, Object> claims = new HashMap<>();
            jwt.getClaims().forEach((k, v) -> claims.put(k, v.as(Object.class)));

            ObjectMapper mapper = new ObjectMapper();
            mapper.enable(SerializationFeature.INDENT_OUTPUT);

            prettyJson = mapper.writeValueAsString(claims);
        }

        req.setAttribute("idToken", idToken);
        req.setAttribute("pretty", prettyJson);

        req.getRequestDispatcher("/WEB-INF/jsp/home.jsp").forward(req, res);
    }
}
