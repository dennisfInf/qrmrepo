// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  fido: {
    attestation: "none",               //"direct" | "enterprise" | "indirect" | "none"
    algorithm: -7,
    authenticatorSelection : {
      authenticatorAttachment: "cross-platform",       // "cross-platform" | "platform"
      requireResidentKey: false,                       // true | false
      residentKey: "preferred",                        // "discouraged" | "preferred" | "required"
      userVerification: "preferred",        // "discouraged" | "preferred" | "required"
    },
    excludeCredentials:[],
    extensions: {
      appid: undefined,
      appidExclude: undefined,
      credProps: false,
      uvm: false,
    },
    pubKeyCredParams : [],
    rp : {
      domain: "localhost",
      name: "elon mask"
    },
    timeout: 60000,
  },
};


