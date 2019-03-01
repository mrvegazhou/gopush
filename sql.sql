CREATE TABLE "public"."t_device" (
"id" int8 DEFAULT nextval('t_device_id_seq'::regclass) NOT NULL,
"user_id" int8 DEFAULT '0'::bigint NOT NULL,
"token" varchar(40) COLLATE "default" NOT NULL,
"type" int2 NOT NULL,
"brand" varchar(20) COLLATE "default" NOT NULL,
"model" varchar(20) COLLATE "default" NOT NULL,
"system_version" varchar(10) COLLATE "default" NOT NULL,
"app_version" varchar(10) COLLATE "default" NOT NULL,
"status" int2 NOT NULL,
"create_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"update_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE)
;

ALTER TABLE "public"."t_device" OWNER TO "dbuser";

COMMENT ON COLUMN "public"."t_device"."id" IS '设备id';

COMMENT ON COLUMN "public"."t_device"."user_id" IS '账户id';

COMMENT ON COLUMN "public"."t_device"."token" IS '设备登录的token';

COMMENT ON COLUMN "public"."t_device"."type" IS '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web';

COMMENT ON COLUMN "public"."t_device"."brand" IS '手机厂商';

COMMENT ON COLUMN "public"."t_device"."model" IS '机型';

COMMENT ON COLUMN "public"."t_device"."system_version" IS '系统版本';

COMMENT ON COLUMN "public"."t_device"."status" IS '在线状态，0：离线；1：在线';

COMMENT ON COLUMN "public"."t_device"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."t_device"."update_time" IS '更新时间';



CREATE TRIGGER "t_name" BEFORE UPDATE ON "public"."t_device"
FOR EACH ROW
EXECUTE PROCEDURE "upd_timestamp"();

create trigger t_name before update on t_device_sync_sequence for each row execute procedure upd_timestamp();
ALTER TABLE t_device_send_sequence ALTER COLUMN create_time SET DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE t_device_send_sequence ALTER COLUMN create_time SET DEFAULT CURRENT_TIMESTAMP;