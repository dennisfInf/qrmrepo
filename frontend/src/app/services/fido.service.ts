import {Injectable} from '@angular/core';
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class FidoService {

  constructor() {
  }

  async createDefaultCredential(challenge: string, displayName: string, userId: string, name: string): Promise<Credential | null> {
    let creationOptions = this.createPublicKeyCredentialCreationOptions(challenge, displayName, userId, name)
    let credential = navigator.credentials.create({publicKey: creationOptions})
    return credential
  }
  async createCredential(challenge: string, displayName: string, userId: string, name: string): Promise<Credential | null> {
    let creationOptions = this.readPublicKeyCredentialCreationOptions(challenge, displayName, userId, name)
    let credential = navigator.credentials.create({publicKey: creationOptions})
    return credential
  }

  private readPublicKeyCredentialCreationOptions(challenge: string, displayName: string, userId: string, name: string): PublicKeyCredentialCreationOptions{
    let publicKeyCredentialCreationOptions: PublicKeyCredentialCreationOptions = {
      attestation: this.readAttestation(),
      authenticatorSelection: this.readAuthenticatorSelection(),
      challenge: Uint8Array.from(challenge, c => c.charCodeAt(0)),
      excludeCredentials: this.readExcludeCredentials(),
      extensions: this.readExtensions(),
      pubKeyCredParams: this.readPubKeyCredParams(),
      rp: this.readRp(),
      timeout: this.readTimeout(),
      user: this.createUser(displayName, userId, name)
    }
    return publicKeyCredentialCreationOptions
  }
  private createPublicKeyCredentialCreationOptions(challenge: string, displayName: string, userId: string, name: string): PublicKeyCredentialCreationOptions {
    let publicKeyCredentialCreationOptions: PublicKeyCredentialCreationOptions = {
      attestation: undefined,
      authenticatorSelection: undefined,
      challenge: Uint8Array.from(challenge, c => c.charCodeAt(0)),
      excludeCredentials: [],
      extensions: undefined,
      pubKeyCredParams: [],
      rp: this.readRp(),
      timeout: this.readTimeout(),
      user: this.createUser(displayName, userId, name)
    }
    return publicKeyCredentialCreationOptions
  }
  private createUser(displayName: string, userId: string, name: string): PublicKeyCredentialUserEntity {
    let publicKeyCredentialUserEntity: PublicKeyCredentialUserEntity = {
      displayName: displayName,
      id: Uint8Array.from(userId, c => c.charCodeAt(0)),
      name: name
    }
    return publicKeyCredentialUserEntity
  }
  private readAuthenticatorSelection(): AuthenticatorSelectionCriteria {
    let authenticatorSelection: AuthenticatorSelectionCriteria = {
      authenticatorAttachment: this.readAuthenticatorAttachment(),
      requireResidentKey: this.readRequireResidentKey(),
      residentKey: this.readResidentKey(),
      userVerification: this.readUserVerification(),
    }
    return authenticatorSelection
  }
  private readTimeout() : number {
    let timeout : number
    timeout = <number> environment.fido.timeout
    return timeout
  }
  private readAuthenticatorAttachment(): AuthenticatorAttachment {
    let authenticatorAttachment: AuthenticatorAttachment
    authenticatorAttachment = <AuthenticatorAttachment>environment.fido.authenticatorSelection.authenticatorAttachment
    return authenticatorAttachment
  }
  private readAttestation(): AttestationConveyancePreference {
    let attestation: AttestationConveyancePreference
    attestation = <AttestationConveyancePreference> environment.fido.attestation
    return attestation
  }
  private readRequireResidentKey(): boolean {
    let requireResidentKey: boolean
    requireResidentKey = environment.fido.authenticatorSelection.requireResidentKey
    return requireResidentKey
  }
  private readResidentKey(): ResidentKeyRequirement {
    let residentKey: ResidentKeyRequirement
    residentKey = <ResidentKeyRequirement>environment.fido.authenticatorSelection.residentKey
    return residentKey
  }
  private readUserVerification(): UserVerificationRequirement {
    let userVerification: UserVerificationRequirement
    userVerification = <UserVerificationRequirement>environment.fido.authenticatorSelection.userVerification
    return userVerification
  }
  private readExtensions(): AuthenticationExtensionsClientInputs {
    let authenticationExtensionsClientInputs: AuthenticationExtensionsClientInputs
    authenticationExtensionsClientInputs = {
      appid: environment.fido.extensions.appid,
      appidExclude: environment.fido.extensions.appidExclude,
      credProps: environment.fido.extensions.credProps,
      uvm: environment.fido.extensions.uvm,
    }
    return authenticationExtensionsClientInputs
  }
  private readExcludeCredentials() : PublicKeyCredentialDescriptor[]{
    let excludeCredentials : PublicKeyCredentialDescriptor[]
    excludeCredentials = <PublicKeyCredentialDescriptor[]> environment.fido.excludeCredentials
    return excludeCredentials
  }
  private readRp() : PublicKeyCredentialRpEntity{
    let rp : PublicKeyCredentialRpEntity
    rp = <PublicKeyCredentialRpEntity> environment.fido.rp
    return rp
  }
  private readPubKeyCredParams() : PublicKeyCredentialParameters[] {
    let pubKeyCredParams : PublicKeyCredentialParameters[]
    pubKeyCredParams = <PublicKeyCredentialParameters[]> environment.fido.pubKeyCredParams
    return pubKeyCredParams
  }
}
