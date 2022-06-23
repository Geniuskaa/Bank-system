//Example of inserting one element
db.user_payments.insertOne({
    login: "boriska",
    payments_transfers: [{
        link_on_icon: "link",
        description: "Purchase on CinemaSite",
        link_on_web_site: "CinemaSite.ru"
    }]
});




//Example of inserting many elements
// db.user_payments.insertMany(
//     [
//         {
//             login: "kilka",
//             payments_transfers: [{
//                 link_on_icon: "link",
//                 description: "Purchase on CinemaSite",
//                 link_on_web_site: "CinemaSite.ru"
//             }]
//         },
//         {
//             login: "filka",
//             payments_transfers: [{
//                 link_on_icon: "link",
//                 description: "Purchase on CinemaSite",
//                 link_on_web_site: "CinemaSite.ru"
//             }]
//
//         }
//     ]);

//Example how to update/add data
// db.user_payments.updateOne({_id: ObjectId('62b418047a292f1ccb6d3616')},
//     {$push: {payments_transfers: {
//                 link_on_icon: "link",
//                 description: "Purchase on WorldOfTanks",
//                 link_on_web_site: "WorldOfTanks.ru"}}});

//Example how to find one document
// db.user_payments.find({login: "boriska"});