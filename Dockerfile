FROM scratch

ADD vault-audit-bridge /bridge

ENTRYPOINT ["/bridge"]
