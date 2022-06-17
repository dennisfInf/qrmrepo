module.exports = {
  "/register/**": {
    secure: false,
    bypass: function(req, res, opts) {
      res.setHeader("Permissions-Policy", "publickey-credentials-get=*");
    }
  }
}
