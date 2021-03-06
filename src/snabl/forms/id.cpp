#include <iostream>

#include "snabl/m.hpp"
#include "snabl/forms/id.hpp"

namespace snabl::forms {
  Id::Id(Pos pos, Sym name): Form(make_shared<const Imp>(pos, name)) {}

  Id::Imp::Imp(Pos pos, Sym name): Form::Imp(pos), name(name) {}

  void Id::Imp::dump(ostream &out) const { out << name; }
  
  optional<Error> Id::Imp::emit(Reg reg, M &m) const {
    optional<Val> v = m.scope->find(name);
    if (!v) { return Error(pos, "Unknown id: ", name); }
    return v->type.imp->methods.emit(*v, reg, pos, m);
  }

  optional<Error> Id::Imp::quote(Reg reg, M &m) const {
    ops::LOAD_SYM(m.emit(), reg, name);
    return nullopt;
  }

  optional<Val> Id::Imp::val(M &m) const { return m.scope->find(name); }
}
