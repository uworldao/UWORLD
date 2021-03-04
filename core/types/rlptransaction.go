package types

import "github.com/uworldao/UWORLD/common/encode/rlp"

type RlpTransaction struct {
	TxHead *TransactionHead
	TxBody []byte
}

func (rt *RlpTransaction) TranslateToTransaction() *Transaction {
	switch rt.TxHead.TxType {
	case NormalTransaction:
		var nt *NormalTransactionBody
		rlp.DecodeBytes(rt.TxBody, &nt)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: nt,
		}
	case ContractTransaction:
		var ct *ContractBody
		rlp.DecodeBytes(rt.TxBody, &ct)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: ct,
		}
	case LoginCandidate:
		var nt *LoginTransactionBody
		rlp.DecodeBytes(rt.TxBody, &nt)
		return &Transaction{
			TxHead: rt.TxHead,
			TxBody: nt,
		}
		/*case LogoutCandidate:
			return &Transaction{
				TxHead: rt.TxHead,
				TxBody: &LogoutTransactionBody{},
			}
		case VoteToCandidate:
			var nt *VoteTransactionBody
			rlp.DecodeBytes(rt.TxBody, &nt)
			return &Transaction{
				TxHead: rt.TxHead,
				TxBody: nt,
			}*/
	}
	return nil
}
