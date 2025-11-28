package com.auth0.example;

import com.auth0.SessionUtils;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class HomeServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse res)
            throws ServletException, IOException {

        String accessToken = (String) SessionUtils.get(req, "accessToken");
        String idToken = (String) SessionUtils.get(req, "idToken");

        // For a real app, you'd parse the ID token or call /userinfo.
        req.setAttribute("accessToken", accessToken);
        req.setAttribute("idToken", idToken);

        req.getRequestDispatcher("/WEB-INF/jsp/home.jsp").forward(req, res);
    }
}
