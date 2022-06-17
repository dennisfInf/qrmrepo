export class Contact {
  name!:string
  address!: string
  imgUrl!: string
}

export const ContactList : Contact[] = [
  {
    name : "John Doe",
    address: "0xb794f5ea0ba39494ce839613fffba742795792689c",
    imgUrl : "https://i.pravatar.cc/100?img=8"
  },
  {
    name : "Volker Racho",
    address: "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
    imgUrl : "https://i.pravatar.cc/100?img=7"
  },
  {
    name : "Rainer Zufall",
    address: "0xbc782de09bf15fe5aa6b10997108a728c4ec7dddx",
    imgUrl : "https://i.pravatar.cc/100?img=6"
  }
]
