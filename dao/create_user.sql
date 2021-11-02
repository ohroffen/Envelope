CREATE DATABASE envelope_rains;

use envelope_rains;

drop table if exists users;

# TODO: 给user_id建索引

# create table users (
#    envelope_id VARCHAR(36)  NOT NULL UNIQUE,
#    user_id VARCHAR(36) NOT NULL UNIQUE,
#    opened BOOLEAN DEFAULT FALSE,
#    value BIGINT(20) DEFAULT 0,
#    snatch_time TIMESTAMP NOT NULL,
#    current_count INTEGER NOT NULL DEFAULT 0,
#
#    PRIMARY KEY (envelope_id)
# ) ENGINE=INNODB DEFAULT CHARSET=utf8;
#
# insert into envelope_rains.users
# (envelope_id, user_id, opened, value, snatch_time, current_count)
# values
#     ('fdge454543', '11111', false, 34, NOW(), 0),
#     ('43534sfdge', '22222', false, 34, NOW(), 0);

create table users (
   envelope_id VARCHAR(36)  NOT NULL UNIQUE,
   user_id VARCHAR(36) NOT NULL,
   opened BOOLEAN DEFAULT FALSE,
   value BIGINT(20) DEFAULT 0,
   snatch_time TIMESTAMP NOT NULL,

   PRIMARY KEY (envelope_id)
) ENGINE=INNODB DEFAULT CHARSET=utf8;

insert into envelope_rains.users
(envelope_id, user_id, opened, value, snatch_time)
values
    ('fdge454543', '11111', false, 34, NOW()),
    ('43534sfdge', '22222', false, 34, NOW());

select * from users;