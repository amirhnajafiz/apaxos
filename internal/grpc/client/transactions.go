package client

type TransactionsDialer struct{}

func (t *TransactionsDialer) NewTransaction() {}

func (t *TransactionsDialer) PrintBalance() {}

func (t *TransactionsDialer) PrintLogs() {}

func (t *TransactionsDialer) PrintDB() {}

func (t *TransactionsDialer) Performance() {}

func (t *TransactionsDialer) AggregatedBalance() {}
