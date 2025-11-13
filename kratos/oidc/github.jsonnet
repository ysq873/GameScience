local claims = std.extVar('claims');
{
  identity: {
    traits: {
      email: if std.objectHas(claims, 'email') then claims.email else if std.objectHas(claims, 'emails') && std.length(claims.emails) > 0 then claims.emails[0].email else '',
      name: if std.objectHas(claims, 'name') then claims.name else if std.objectHas(claims, 'preferred_username') then claims.preferred_username else ''
    }
  }
}
