import Image from "next/image";

export function About() {
  return (
    <div className="bg-off-white p-20">
      <div className="flex flex-col bg-[#C4C4C4] py-20 items-center gap-6">
        <p className="font-[Syne]">Welcome to HeartLink</p>
        <h1 className="text-4xl md:text-5xl font-[Playfair Display SC] font-bold">
          Revolutionizing Low-Cost Telemedicine Solutions
        </h1>
        <p className="font-[Syne]">
          At HeartLink, we are dedicated to providing innovative healthcare
          solutions that empower individuals to take control of their
          well-being.
        </p>
      </div>
      <div className="p-20 pt-40">
        <div className="flex gap-20">
          <div className="flex flex-col gap-5 items-center">
            <Image
              width="96"
              height="96"
              src="https://img.icons8.com/color/96/heart-with-pulse.png"
              alt="heart-with-pulse"
            />
            <h2 className="text-3xl font-[Playfair Display SC] font-bold">
              Our Mission
            </h2>
            <p className="font-[Syne]">
              Heart murmurs affect 30% of children and 10% of adults, often
              signaling serious conditions. Rural Canadians face diagnosis
              barriers due to provider shortages and long travel distances.
              Telehealth can help, but affordable at-home diagnostic tools are
              needed — That’s where HeartLink comes in.
            </p>
          </div>
          <div className="flex flex-col gap-5 items-center">
            <Image
              width="100"
              height="100"
              src="https://img.icons8.com/plasticine/100/stethoscope.png"
              alt="stethoscope"
            />
            <h2 className="text-3xl font-[Playfair Display SC] font-bold">
              Remote Diagnosis
            </h2>
            <p className="font-[Syne]">
              Doctors can review patient data and heart sounds remotely,
              enabling faster diagnoses, treatment monitoring, and seamless
              collaboration. Our platform connects physicians and patients,
              making expert care accessible anywhere.
            </p>
          </div>
          <div className="flex flex-col gap-5 items-center">
            <Image
              src="https://img.icons8.com/color/96/medical-mobile-app.png"
              width="96"
              height="96"
              alt="medical-mobile-app"
            />
            <h2 className="text-3xl font-[Playfair Display SC] font-bold">
              Download the App
            </h2>
            <p className="font-[Syne]">
              Patients can take control of their health with our mobile app,
              accessing their data, diagnoses, and treatment plans in one place.
              The app connects directly to the digital stethoscope, enabling
              remote doctor communication and effortless care management.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
