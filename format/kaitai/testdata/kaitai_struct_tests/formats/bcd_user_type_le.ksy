# https://github.com/kaitai-io/kaitai_struct/issues/78
meta:
  id:      bcd_user_type_le
  endian:  le

seq:
  - id:       ltr
    size:     4
    type:     ltr_obj
  - id:       rtl
    size:     4
    type:     rtl_obj
  - id:       leading_zero_ltr
    size:     4
    type:     leading_zero_ltr_obj

types:
  ltr_obj:
    seq:
      - id:   b1
        type: u1
      - id:   b2
        type: u1
      - id:   b3
        type: u1
      - id:   b4
        type: u1

    instances:
      digit1:
        value: (b4 & 0xF0) >> 4
      digit2:
        value: (b4 & 0x0F)
      digit3:
        value: (b3 & 0xF0) >> 4
      digit4:
        value: (b3 & 0x0F)
      digit5:
        value: (b2 & 0xF0) >> 4
      digit6:
        value: (b2 & 0x0F)
      digit7:
        value: (b1 & 0xF0) >> 4
      digit8:
        value: (b1 & 0x0F)
      as_int:
        value:  digit8 *        1 +
                digit7 *       10 +
                digit6 *      100 +
                digit5 *     1000 +
                digit4 *    10000 +
                digit3 *   100000 +
                digit2 *  1000000 +
                digit1 * 10000000
      as_str:
        value:  digit1.to_s + digit2.to_s + digit3.to_s + digit4.to_s + digit5.to_s + digit6.to_s + digit7.to_s + digit8.to_s

  rtl_obj:
    seq:
      - id:   b1
        type: u1
      - id:   b2
        type: u1
      - id:   b3
        type: u1
      - id:   b4
        type: u1

    instances:
      digit1:
        value: (b4 & 0xF0) >> 4
      digit2:
        value: (b4 & 0x0F)
      digit3:
        value: (b3 & 0xF0) >> 4
      digit4:
        value: (b3 & 0x0F)
      digit5:
        value: (b2 & 0xF0) >> 4
      digit6:
        value: (b2 & 0x0F)
      digit7:
        value: (b1 & 0xF0) >> 4
      digit8:
        value: (b1 & 0x0F)
      as_int:
        value:  digit1 *        1 +
                digit2 *       10 +
                digit3 *      100 +
                digit4 *     1000 +
                digit5 *    10000 +
                digit6 *   100000 +
                digit7 *  1000000 +
                digit8 * 10000000
      as_str:
        value:  digit8.to_s + digit7.to_s + digit6.to_s + digit5.to_s + digit4.to_s + digit3.to_s + digit2.to_s + digit1.to_s

  leading_zero_ltr_obj:
    seq:
      - id:   b1
        type: u1
      - id:   b2
        type: u1
      - id:   b3
        type: u1
      - id:   b4
        type: u1

    instances:
      digit1:
        value: (b4 & 0xF0) >> 4
      digit2:
        value: (b4 & 0x0F)
      digit3:
        value: (b3 & 0xF0) >> 4
      digit4:
        value: (b3 & 0x0F)
      digit5:
        value: (b2 & 0xF0) >> 4
      digit6:
        value: (b2 & 0x0F)
      digit7:
        value: (b1 & 0xF0) >> 4
      digit8:
        value: (b1 & 0x0F)
      as_int:
        value:  digit8 *        1 +
                digit7 *       10 +
                digit6 *      100 +
                digit5 *     1000 +
                digit4 *    10000 +
                digit3 *   100000 +
                digit2 *  1000000 +
                digit1 * 10000000
      as_str:
        value:  digit1.to_s + digit2.to_s + digit3.to_s + digit4.to_s + digit5.to_s + digit6.to_s + digit7.to_s + digit8.to_s
