// Project: PA193_mnemonic_aller
// Maintainers UCO: 408788, 497391, 497577
// Description: Mnemonic API

package mnemonic

import (
  "fmt"
  "crypto/sha512"
  "encoding/hex"
)

func EntropyToPhraseAndSeed() {
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
}

func PhraseToEntropyAndSeed(phrase string) (entropy string,seed string){
  // TODO
  fmt.Println("TODO entropy_to_mnemonic")
  // Trasnformation of phrase to entropy and seed
  seed = phrase     // For test, must be deleted and replaced
  return entropy, seed
}

func VerifyPhraseAndSeed(phrase_to_verify,seed_to_verify string) int{
  // TODO
  // fmt.Println("TODO entropy_to_mnemonic")
  var seed string
  _, seed = PhraseToEntropyAndSeed(phrase_to_verify)
  if seed_to_verify==seed{
    fmt.Println("The phrase and the seed correspond !!")
    return 0
  } else {
    fmt.Println("The phrase and the seed do NOT correspond !!")
    return -1
  }
  return -1
}


func PBKDF2_SHA512_F(password, salt []byte, count, l_counter int)([]byte, int){

  U_1_to_c := make([][]byte,count)

  var INT_i_l [4]byte
  INT_i_l[0] = byte(l_counter >> 24)
  INT_i_l[1] = byte(l_counter >> 16)
  INT_i_l[2] = byte(l_counter >> 8)
  INT_i_l[3] = byte(l_counter)


  U_1_to_c[0] = make([]byte,64)
  sha_512 := sha512.New()
  sha_512.Write(password)
  // fmt.Println("password : ",hex.EncodeToString(password))
  sha_512.Write(append(salt,INT_i_l[:4]...))
  // fmt.Println("1 salt : ",hex.EncodeToString(append(salt,INT_i_l[:4]...)))
  U_1_to_c[0] = sha_512.Sum(nil)
  sha_512.Reset()
  // fmt.Println("U_1_to_c[0] : ",hex.EncodeToString(U_1_to_c[0]))

  for i:=1 ; i<count ; i++ {
    sha_512 = sha512.New()
    U_1_to_c[i] = make([]byte,64)
    sha_512.Write(password)
    sha_512.Write(U_1_to_c[i-1])
    U_1_to_c[i] = sha_512.Sum(nil)
    sha_512.Reset()
  }


  output := make([]byte,64)
  output = U_1_to_c[0]
  for i:=1; i<count ; i++ {
    for j := range U_1_to_c[i]{
      output[j] ^= U_1_to_c[i][j]
    }

  }

  return output, 0
}


func PBKDF2_SHA512(password, salt []byte, count, output_len int)([]byte, int){
  if(output_len != 64){ /* Length of SHA-512 */
    return nil,-1
  } else {
    fmt.Println(hex.EncodeToString(password),salt)
    hLen := 64  /* Length of SHA-512 */
    l := output_len / hLen  /* Should be equal to 1 !*/
    // r := output_len -(l-1)*hLen     /* Should be equal to output_len, so 64 bytes */

    var T_1_to_l [][]byte
    T_1_to_l = make([][]byte,output_len)

    for i:=0; i<l ; i++ {
      T_1_to_l[i] = make([]byte,output_len)
      T_1_to_l[i],_ = PBKDF2_SHA512_F(password, salt,count,i+1)

    }

    output := make([]byte,output_len)
    output = T_1_to_l[0]
    for i:=1 ; i<l ; i++ {
      output = append(output,T_1_to_l[i]...)
    }
    return output,0
  }
}



func PhraseToSeed(phrase,passphrase string)(seed []byte,err int){
  seed , err = PBKDF2_SHA512([]byte(phrase),[]byte("mnemonic"+passphrase),2048,64)
  return seed, err
}
