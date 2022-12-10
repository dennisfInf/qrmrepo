
export const environment = {
  production: true,
  fido: {
    attestation: "none",               //"direct" | "enterprise" | "indirect" | "none"
    algorithm: -7,
    authenticatorSelection: {
      authenticatorAttachment: "platform",       // "cross-platform" | "platform"
      requireResidentKey: false,                       // true | false
      residentKey: "discouraged",                        // "discouraged" | "preferred" | "required"
      userVerification: "required",        // "discouraged" | "preferred" | "required"
    },
    excludeCredentials: [],
    extensions: {
      appid: undefined,
      appidExclude: undefined,
      credProps: false,
      uvm: false,
    },
    pubKeyCredParams: [{ type: "public-key", alg: -7 }],
    rp : {
      domain: "elonwallet.io",
      name: "elonwallet"
    },
    timeout: 60000,
  },
  routes : {
    // @ts-ignore
    authenticationService: "https://elonwallet.io/api"
  }
};
