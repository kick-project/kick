SELECT base.base AS dir, file.file AS path, label.label AS label FROM file
	LEFT JOIN base ON base.id = file.base_id
	LEFT JOIN file_label ON file.id = file_label.file_id
	LEFT JOIN label ON file_label.label_id = label.id