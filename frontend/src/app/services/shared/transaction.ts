import {ContactList} from "./contact";

export class Transaction {
  from!:string
  to!: string
  contactName!: string
  amount!: string
  date!: string
}

export const TransactionList : Transaction[] = [
  {
    from: ContactList[1].address,
    to : "0x8B4A4EC8303C24685D276C60A4C74286C6aC4D",
    contactName: ContactList[1].name,
    amount: "0.5",
    date : "01/01/2022"
  },
  {
    from: "0x8B4A4EC8303C24685D276C60A4C74286C6aC4D",
    to : ContactList[2].address,
    contactName: ContactList[2].name,
    amount: "1.3",
    date : "01/01/2022"

  },
  {
    from: ContactList[0].address,
    to : "0x8B4A4EC8303C24685D276C60A4C74286C6aC4D",
    contactName: ContactList[0].name,
    amount: "0.1",
    date : "01/01/2022"

  },
  {
    from: ContactList[1].address,
    to :  "0x8B4A4EC8303C24685D276C60A4C74286C6aC4D",
    contactName: ContactList[1].name,
    amount: "0.4",
    date : "01/01/2022"

  },

  {
    from: "0x8B4A4EC8303C24685D276C60A4C74286C6aC4D",
    to : ContactList[2].address,
    contactName: ContactList[2].name,
    amount: "0.9",
    date : "01/01/2022"
  }
]
