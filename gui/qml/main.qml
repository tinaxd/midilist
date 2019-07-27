import QtQuick 2.13
import QtQuick.Window 2.13
import QtQuick.Controls 2.13
import QtQuick.Controls 1.4
import QtQuick.Dialogs 1.0
import CustomQmlTypes 1.0

Window {
    visible: true
    width: 640
    height: 480
    title: qsTr("Midilist")

    property BridgeTemplate template: BridgeTemplate{}

    Row {
        id: row
        spacing: 0
        anchors.fill: parent

        ToolBar {
            id: menuBar
            anchors.left: parent.left
            anchors.leftMargin: 0
            anchors.right: parent.right
            anchors.rightMargin: 0

            Button {
                id: readBtn
                text: qsTr("Load")
                onClicked: {
                    midiReadDialog.visible = true
                }
            }

            Button {
                id: quitBtn
                text: qsTr("Quit")

                onClicked: {
                    root.template.quitBtnClicked()
                }
            }
        }

        Flow {
            id: element
            anchors.bottom: bottomBar.top
            anchors.bottomMargin: 0
            anchors.top: menuBar.bottom
            anchors.topMargin: 0
            anchors.right: parent.left
            anchors.rightMargin: -200
            anchors.left: parent.left
            anchors.leftMargin: 0

            TableView {
                id: eventView
                anchors.right: parent.left
                anchors.rightMargin: -400
                anchors.left: parent.left
                anchors.leftMargin: 0
                anchors.top: parent.top
                anchors.topMargin: 0
                anchors.bottom: parent.bottom
                anchors.bottomMargin: 0

                model: CustomTableModel{}
                TableViewColumn {
                    role: "Time"
                    title: role
                    width: 125
                }
                TableViewColumn {
                    role: "Event"
                    title: role
                    width: 225
                }
                TableViewColumn {
                    role: "Value"
                    title: role
                    width: 50
                }
            }

            ListModel {
                id: eventModel
                ListElement {
                    time: 0
                    event: "RPN LSB"
                    value: 0
                }
            }
        }

        ToolBar {
            id: bottomBar
            anchors.left: parent.left
            anchors.leftMargin: 0
            anchors.right: parent.right
            anchors.rightMargin: 0
            anchors.bottom: parent.bottom
            anchors.bottomMargin: 0
        }

    }

    FileDialog {
        id: midiReadDialog
        title: "Please choose a midi file"
        foler: shortcuts.home
        onAccepted: {
        }
        onRejected: {
        }
        visible: true
    }
}

