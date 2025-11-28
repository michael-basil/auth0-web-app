// web/static/js/user.js
$(document).ready(function () {
  $(".btn-logout").click(function () {
    Cookies.remove("auth-session");
  });
});
