from __future__ import annotations

from PyQt6.QtCore import QDate
from PyQt6.QtGui import QStandardItem, QStandardItemModel
from PyQt6.QtWidgets import (
    QComboBox, QDateEdit, QDialog, QDialogButtonBox, QDoubleSpinBox,
    QFormLayout, QHBoxLayout, QLineEdit, QMessageBox, QPushButton,
    QSpinBox, QTableView, QTextEdit, QVBoxLayout, QWidget,
)

from db_manager import DatabaseError, DatabaseManager
from table_specs import FieldSpec, FieldType, TableSpec


class RecordDialog(QDialog):
    """Форма добавления / редактирования одной записи, построенная по FieldSpec."""

    def __init__(self, spec: TableSpec, db: DatabaseManager, values: dict | None = None, parent=None):
        super().__init__(parent)
        self.spec = spec
        self.db = db
        self.values = values or {}
        self.editing = values is not None
        self.setWindowTitle(("Изменить" if self.editing else "Добавить") + f" — {spec.title}")

        self.widgets: dict[str, object] = {}
        layout = QFormLayout(self)

        for f in spec.fields:
            widget = self._build_widget(f)
            if self.editing and f.name == spec.pk:
                # первичный ключ существующей записи не меняем во избежание
                # рассинхронизации со связанными таблицами
                widget.setEnabled(False)
            self.widgets[f.name] = widget
            layout.addRow(f.label + ":", widget)

        buttons = QDialogButtonBox(
            QDialogButtonBox.StandardButton.Ok | QDialogButtonBox.StandardButton.Cancel
        )
        buttons.accepted.connect(self._on_accept)
        buttons.rejected.connect(self.reject)
        layout.addRow(buttons)

    def _build_widget(self, f: FieldSpec):
        current = self.values.get(f.name)

        if f.field_type == FieldType.LOOKUP:
            combo = QComboBox()
            try:
                options = self.db.fetch_lookup(f.lookup_table, f.lookup_key, f.lookup_label)
            except DatabaseError as exc:
                QMessageBox.critical(self, "Ошибка", str(exc))
                options = []
            for key, label in options:
                combo.addItem(f"{label} ({key})", key)
            if current is not None:
                idx = combo.findData(current)
                if idx >= 0:
                    combo.setCurrentIndex(idx)
            combo.setEnabled(f.editable)
            return combo

        if f.field_type == FieldType.TEXTAREA:
            edit = QTextEdit()
            edit.setPlainText("" if current is None else str(current))
            edit.setEnabled(f.editable)
            edit.setFixedHeight(70)
            return edit

        if f.field_type == FieldType.INT:
            spin = QSpinBox()
            spin.setRange(-2_147_483_648, 2_147_483_647)
            spin.setValue(int(current) if current is not None else 0)
            spin.setEnabled(f.editable)
            return spin

        if f.field_type == FieldType.DECIMAL:
            spin = QDoubleSpinBox()
            spin.setRange(0, 999_999_999)
            spin.setDecimals(2)
            spin.setValue(float(current) if current is not None else 0.0)
            spin.setEnabled(f.editable)
            return spin

        if f.field_type == FieldType.DATE:
            date_edit = QDateEdit()
            date_edit.setCalendarPopup(True)
            date_edit.setDisplayFormat("yyyy-MM-dd")
            if current:
                date_edit.setDate(QDate(current.year, current.month, current.day))
            else:
                date_edit.setDate(QDate.currentDate())
            date_edit.setEnabled(f.editable)
            return date_edit

        edit = QLineEdit()
        edit.setText("" if current is None else str(current))
        edit.setEnabled(f.editable)
        return edit

    def _on_accept(self) -> None:
        for f in self.spec.fields:
            if not f.required:
                continue
            widget = self.widgets[f.name]
            text = None
            if f.field_type == FieldType.TEXT:
                text = widget.text().strip()
            elif f.field_type == FieldType.TEXTAREA:
                text = widget.toPlainText().strip()
            if text is not None and not text:
                QMessageBox.warning(self, "Проверка данных", f"Поле «{f.label}» обязательно для заполнения.")
                return
        self.accept()

    def get_values(self) -> dict:
        result = {}
        for f in self.spec.fields:
            widget = self.widgets[f.name]
            if f.field_type == FieldType.LOOKUP:
                result[f.name] = widget.currentData()
            elif f.field_type == FieldType.TEXTAREA:
                result[f.name] = widget.toPlainText().strip()
            elif f.field_type == FieldType.INT:
                result[f.name] = widget.value()
            elif f.field_type == FieldType.DECIMAL:
                result[f.name] = widget.value()
            elif f.field_type == FieldType.DATE:
                qd = widget.date()
                result[f.name] = f"{qd.year():04d}-{qd.month():02d}-{qd.day():02d}"
            else:
                result[f.name] = widget.text().strip()
        return result


class CrudTab(QWidget):
    """Вкладка одной таблицы: список записей + добавление/изменение/удаление."""

    def __init__(self, spec: TableSpec, db: DatabaseManager, parent=None):
        super().__init__(parent)
        self.spec = spec
        self.db = db
        self.rows: list[dict] = []

        self.table_view = QTableView()
        self.table_view.setSelectionBehavior(QTableView.SelectionBehavior.SelectRows)
        self.table_view.setSelectionMode(QTableView.SelectionMode.SingleSelection)
        self.table_view.setEditTriggers(QTableView.EditTrigger.NoEditTriggers)
        self.model = QStandardItemModel()
        self.table_view.setModel(self.model)

        add_btn = QPushButton("Добавить")
        edit_btn = QPushButton("Изменить")
        delete_btn = QPushButton("Удалить")
        refresh_btn = QPushButton("Обновить")
        add_btn.clicked.connect(self.on_add)
        edit_btn.clicked.connect(self.on_edit)
        delete_btn.clicked.connect(self.on_delete)
        refresh_btn.clicked.connect(self.refresh)

        toolbar = QHBoxLayout()
        for b in (add_btn, edit_btn, delete_btn, refresh_btn):
            toolbar.addWidget(b)
        toolbar.addStretch()

        layout = QVBoxLayout(self)
        layout.addLayout(toolbar)
        layout.addWidget(self.table_view)

        self.refresh()

    def refresh(self) -> None:
        try:
            self.rows = self.db.fetch_all(self.spec.table, order_by=self.spec.pk)
        except DatabaseError as exc:
            QMessageBox.critical(self, "Ошибка", str(exc))
            self.rows = []
        self._populate_model()

    def _populate_model(self) -> None:
        self.model.clear()
        if not self.rows:
            self.model.setHorizontalHeaderLabels([f.label for f in self.spec.fields])
            return
        headers = list(self.rows[0].keys())
        self.model.setHorizontalHeaderLabels(headers)
        for row in self.rows:
            items = [QStandardItem("" if v is None else str(v)) for v in row.values()]
            for it in items:
                it.setEditable(False)
            self.model.appendRow(items)
        self.table_view.resizeColumnsToContents()

    def _selected_row(self) -> dict | None:
        sel = self.table_view.selectionModel().selectedRows()
        if not sel:
            return None
        return self.rows[sel[0].row()]

    def on_add(self) -> None:
        dialog = RecordDialog(self.spec, self.db, values=None, parent=self)
        if dialog.exec() == QDialog.DialogCode.Accepted:
            values = dialog.get_values()
            try:
                self.db.insert(self.spec.table, values)
            except DatabaseError as exc:
                QMessageBox.critical(self, "Ошибка добавления", str(exc))
                return
            self.refresh()

    def on_edit(self) -> None:
        row = self._selected_row()
        if row is None:
            QMessageBox.information(self, "Изменение", "Выберите запись для изменения.")
            return
        dialog = RecordDialog(self.spec, self.db, values=row, parent=self)
        if dialog.exec() == QDialog.DialogCode.Accepted:
            values = dialog.get_values()
            values.pop(self.spec.pk, None)  # pk не меняем при редактировании
            try:
                self.db.update(self.spec.table, self.spec.pk, row[self.spec.pk], values)
            except DatabaseError as exc:
                QMessageBox.critical(self, "Ошибка изменения", str(exc))
                return
            self.refresh()

    def on_delete(self) -> None:
        row = self._selected_row()
        if row is None:
            QMessageBox.information(self, "Удаление", "Выберите запись для удаления.")
            return
        confirm = QMessageBox.question(
            self, "Удаление", "Удалить выбранную запись?",
            QMessageBox.StandardButton.Yes | QMessageBox.StandardButton.No,
        )
        if confirm != QMessageBox.StandardButton.Yes:
            return
        try:
            self.db.delete(self.spec.table, self.spec.pk, row[self.spec.pk])
        except DatabaseError as exc:
            QMessageBox.critical(self, "Ошибка удаления", str(exc))
            return
        self.refresh()
