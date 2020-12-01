INSERT INTO configuration (key, tenant, type, name, date_from, date_to, value) VALUES ('smtp.server.port', '', 'Int', 'SMTP port', '2020-01-01', '2020-12-31', '25');
INSERT INTO configuration (key, tenant, type, name, date_from, date_to, value) VALUES ('smtp.server.port', '', 'Int', 'SMTP port', '2021-01-01', '2999-12-31', '25');
INSERT INTO configuration (key, tenant, type, name, date_from, date_to, value) VALUES ('smtp.server.port', '', 'Int', 'SMTP port', '2019-01-01', '2019-12-31', '25');
INSERT INTO configuration (key, tenant, type, name, date_from, date_to, value) VALUES ('smtp.server.host', '', 'String', 'SMTP Host name or IP', '2020-01-01', '2999-12-31', '127.0.0.1');

INSERT INTO dictionary_metadata (type, tenant, content) VALUES ('AbsenceType', '', '{"$id": "https://alapierre.io/dictionary.schema.json", "type": "object", "title": "DictionaryAbsenceType", "$schema": "http://json-schema.org/draft-07/schema#", "required": ["onlyOnBeginOrEnd", "needDeliveryDateConfirmation", "needConfirmationDocumentNumber"], "properties": {"onlyOnBeginOrEnd": {"type": "boolean", "default": false, "description": "Absence can only start on beginning or finish on end of work day"}, "needDeliveryDateConfirmation": {"type": "boolean", "default": false, "description": "Is proof of absence delivery date required - should field be visible on form"}, "needConfirmationDocumentNumber": {"type": "boolean", "default": false, "description": "Is absence confirmation document number needed"}}}');

INSERT INTO dictionary.dictionary (key, type, group_id, tenant, content, name, parent_key, lp) VALUES ('uopd', 'AbsenceType', null, '', '{"onlyOnBeginOrEnd": true, "needDeliveryDateConfirmation": true, "needConfirmationDocumentNumber": false}', 'Urlop Opieka nad dzieckiem', null, '');
INSERT INTO dictionary.dictionary (key, type, group_id, tenant, content, name, parent_key, lp) VALUES ('uopdd', 'AbsenceType', null, '', '{"label": "Opieka rozliczana dziennie"}', 'Child 1', 'uopd', '');
INSERT INTO dictionary.dictionary (key, type, group_id, tenant, content, name, parent_key, lp) VALUES ('uopdh', 'AbsenceType', null, '', '{"label": "Opieka godzinowa"}', 'Child 2', 'uopd', '');
