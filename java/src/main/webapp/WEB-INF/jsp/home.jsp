<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>Auth0 Java Example</title>
</head>
<body>
<h1>Welcome!</h1>

<p>Access Token:</p>
<pre>${accessToken}</pre>

<p>ID Token:</p>
<pre>${idToken}</pre>
<pre>${pretty}</pre>

<p><a href="${pageContext.request.contextPath}/logout">Logout</a></p>
</body>
</html>
