migrateup:
		goose up

new_migrate:
		goose create ${name} sql


.PHONY: migrateup new_migrate