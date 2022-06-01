export const environment = {
  production: true,
  fido: {
    attestation: "none",               //"direct" | "enterprise" | "indirect" | "none"
    algorithm: -7,
    authenticatorSelection : {
      authenticatorAttachment: "platform",       // "cross-platform" | "platform"
      requireResidentKey: false,                       // true | false
      residentKey: "discouraged",                        // "discouraged" | "preferred" | "required"
      userVerification: "required",        // "discouraged" | "preferred" | "required"
    },
    excludeCredentials:[],
    extensions: {
      appid: undefined,
      appidExclude: undefined,
      credProps: false,
      uvm: false,
    },
    pubKeyCredParams : [ {type: "public-key", alg: -7}],
    rp : {
      domain: "localhost",
      name: "elon mask"
    },
    timeout: 60000,
  },
  routes : {
    authenticationService: "http://localhost:8080"
  }
};
