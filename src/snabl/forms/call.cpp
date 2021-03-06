#include <iostream>

#include "snabl/m.hpp"
#include "snabl/macro.hpp"
#include "snabl/forms/call.hpp"
#include "snabl/forms/id.hpp"

namespace snabl::forms {
  Call::Call(Pos pos, Form target, const deque<Form> &args): Form(make_shared<const Imp>(pos, target, args)) {}

  Call::Imp::Imp(Pos pos, Form target, const deque<Form> &args):
    Form::Imp(pos), target(target), args(args.begin(), args.end()) {}

  void Call::Imp::dump(ostream &out) const {
    out << '(';
    target.dump(out);

    for (Form f: args) {
      out << ' ';
      f.dump(out);
    }
    
    out << ')';
  }
  
  optional<Error> Call::Imp::emit(Reg reg, M &m) const {
    Sym target_id = target.as<Id>().name;
    optional<Val> v(m.scope->find(target_id));
    if (!v) { return Error(pos, "Unknown call target: ", target_id); }
    
    if (v->type == m.abc_lib->macro_type) { return v->as<snabl::Macro *>()->emit(args, reg, pos, m); }

    ops::STATE_BEG(m.emit(), m.emit_pc);

    for (int i = 0; i < args.size(); i++) {
      if (auto err = args[i].emit(i+1, m); err) { return err; }
    }

    if (v->type == m.abc_lib->reg_type) { ops::CALL(m.emit(), v->as<Reg>(), reg); }
    else {
      if (v->type != m.abc_lib->fun_type) { return Error(pos, "Invalid call target: ", *v); }
      Fun *fun = v->as<Fun *>();
      
      if (reinterpret_cast<Op>(fun) <= CALLI1_TARGET_MAX) {
	ops::CALLI1(m.emit(), reg, fun);
      } else {
	Reg fun_reg = m.scope->reg_count++;
	ops::LOAD_FUN(m.emit(2), fun_reg, fun);
	ops::CALL(m.emit(), fun_reg, reg);
      }
    }
    
    return nullopt;
  }
}
