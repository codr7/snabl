#ifndef SNABL_STATE_HPP
#define SNABL_STATE_HPP

#include <array>

#include "snabl/val.hpp"

#define OP_REG_BITS 6

namespace snabl {
  using namespace std;
    
  struct State {
    static const int REG_COUNT = 1 << OP_REG_BITS;
    
    State *outer;
    array<optional<Val>, REG_COUNT> _regs;
    int ref_count;

    State(State *outer): outer(outer), ref_count(1) {
      if (outer) { outer->ref_count++; }
    }
   
    optional<Val> &find(Reg reg) {
      optional<Val> &v = _regs[reg];
      return (!v && outer) ? outer->find(reg) : v;
    }

    Val &get(Reg reg) { return *find(reg); }
    
    void set(Reg reg, const Val &val) { _regs[reg] = val; }
    void set(Reg reg, Val &&val) { _regs[reg] = move(val); }
  };
}

#endif
