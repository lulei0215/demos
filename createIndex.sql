CREATE INDEX idx_socket_id ON task (socket_id);

CREATE INDEX idx_oaid ON device_list(oaid);

ALTER TABLE `oceanengine` ADD INDEX `idx_adid_cid` (`adid`, `cid`);

CREATE INDEX idx_task_status_create_time ON task(status, create_time);

CREATE INDEX idx_task_status ON task(status);

CREATE INDEX idx_imei ON device_list (imei);

CREATE INDEX idx_create_time ON device_list(create_time);

CREATE INDEX idx_oceanengine_account_id ON oceanengine(account_id);

CREATE INDEX idx_chat_record_comb ON chat_record (user_id, is_delete, type, modelId, create_time);

CREATE INDEX idx_combined1 ON chat_record(create_time, role, modelId, user_id);

CREATE INDEX idx_user_account_token ON user_account (user_token);

CREATE INDEX idx_register_invite_code ON user_account(register_invite_code);

CREATE INDEX idx_account ON user_account (account);

CREATE INDEX idx_create_time_action_seq ON user_action (create_time, action_seq, device_id);
CREATE INDEX idx_mac_ua ON device_list(mac, ua);
CREATE INDEX idx_user_action ON user_action (user_id, action_event, action_seq, create_time);

CREATE INDEX idx_user_action_create_time ON user_action(create_time);
CREATE INDEX idx_user_account_vip_time ON user_account(vip_time);


ALTER TABLE `aiworld`.`device_list` ADD INDEX `idx_androidid` (`android_id`)
ALTER TABLE `aiworld`.`device_list` ADD INDEX `idx_oaid` (`oaid`)

ALTER TABLE `aiworld`.`device_list` ADD INDEX `idx_operip_mac_ua` (`oper_ip`, `mac`, `ua`)




DELETE FROM device_list WHERE JSON_EXTRACT(record_info, '$.activationTime') = 0 and TIMESTAMPDIFF(HOUR, update_time, NOW()) > 48;


SELECT COUNT(*) FROM device_list WHERE JSON_EXTRACT(record_info, '$.activationTime') != 0


DELIMITER //

CREATE TRIGGER after_transaction_update
AFTER UPDATE ON transaction
FOR EACH ROW
BEGIN
    IF NEW.amount > 0 AND NEW.expiration_time IS NOT NULL THEN
        INSERT INTO user_invitation_reward
            (create_time, update_time, create_by, update_by, remark, version, user_id, reward_day, start_time, end_time, expiration_time, name, type)
        VALUES
            (NEW.create_time, CURRENT_TIMESTAMP(6), '', '', NEW.vip_id, 1, NEW.user_id, 0, NEW.create_time, NEW.expiration_time, NEW.expiration_time, '', 3);
    END IF;
END;

//
DELIMITER ;


select  COUNT(*) FROM  device_list where  JSON_EXTRACT(record_info, '$.activationTime') != 0  and JSON_EXTRACT(record_info, '$.registerTime') = 0 and os = 0 and create_time > '2024-01-15 00:00:00'

select  COUNT(*) FROM  device_list where  JSON_EXTRACT(record_info, '$.activationTime') != 0 and os = 0 and create_time > '2024-01-15 00:00:00'
